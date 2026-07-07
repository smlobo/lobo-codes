//
// Created by Sheldon Lobo on 7/6/26.
//

#include "edge-weighted-graph.h"

#include <cassert>
#include <iostream>

#include "main.h"
#include "vertex.h"

EdgeWeightedGraph::EdgeWeightedGraph(unsigned nVertices, Context& ctx) : nextId(0) {

    // Create random vertices
    while (vertices.size() < nVertices) {
        int x = ctx.uniformXLocation(ctx.re);
        int y = ctx.uniformYLocation(ctx.re);

        bool tooClose = false;
        for (const auto& idVertexPair : vertices) {
            if (idVertexPair.second.tooClose(x, y)) {
                std::cout << "Skipping random vertex: " << x << ", " << y << "; near: " << idVertexPair.second << std::endl;
                tooClose = true;
                break;
            }
        }
        if (!tooClose) {
            auto [idVertexPair, inserted] = vertices.emplace(nextId, Vertex(x, y, nextId));
            assert(inserted);
            nextId++;
            std::cout << "Created random vertex: " << idVertexPair->second << std::endl;
        }
    }

    // Create edges for a complete graph
    std::vector<Edge> completeEdges;
    for (auto it1 = vertices.begin(); it1 != vertices.end(); ++it1) {
        auto it2 = std::next(it1);
        for (; it2 != vertices.end(); ++it2) {
            Vertex* v1 = &(it1->second);
            Vertex* v2 = &(it2->second);
            Vertex *from, *to = nullptr;
            if (v1->euclideanDistance < v2->euclideanDistance) {
                from = v1;
                to = v2;
            } else {
                from = v2;
                to = v1;
            }
            completeEdges.emplace_back(from, to);
        }
    }
    std::sort(completeEdges.begin(), completeEdges.end(), EdgeWeightComparator());
    // Add all non-intersecting edges to the graph
    for (Edge& completeEdge : completeEdges) {
        std::cout << "Consider edge: " << completeEdge << std::endl;
        bool intersects = false;
        for (auto& existingEdge : edges) {
            if (completeEdge.intersects(&existingEdge)) {
                intersects = true;
                std::cout << "  Intersects edge: " << existingEdge << std::endl;
                break;
            }
        }
        if (!intersects) {
            Vertex *from = completeEdge.from;
            Vertex *to = completeEdge.to;
            auto [it, inserted] = edges.emplace(from, to);
            assert(inserted);
            const Edge* edgePtr = &(*it);
            from->edges.emplace_back(edgePtr);
            to->edges.emplace_back(edgePtr);
            std::cout << "  Created edge: " << *edgePtr << std::endl;
        }
    }
}

void EdgeWeightedGraph::removeVertex(unsigned id) {
    // Get pointer
    auto idVertexPair = vertices.find(id);
    assert(idVertexPair != vertices.end());
    Vertex* vertex = &idVertexPair->second;

    // List of vertices that could have new edges between each other
    std::vector<Vertex*> affectedVertexes;

    // Remove edges
    for (unsigned i = 0; i < vertex->edges.size(); i++) {
        const Edge *edge = vertex->edges[i];
        std::cout << "  Removing edge: " << *edge << std::endl;
        if (edge->from == vertex) {
            affectedVertexes.push_back(edge->to);
            edge->to->removeEdge(edge);
        } else {
            assert(edge->to == vertex);
            affectedVertexes.push_back(edge->from);
            edge->from->removeEdge(edge);
        }
        edges.erase(*edge);
    }
    vertex->edges.clear();

    vertices.erase(idVertexPair);

    std::cout << "Affected vertexes: " << affectedVertexes.size() << std::endl;
    // Create edges for a complete graph
    std::vector<Edge> completeEdges;
    for (std::vector<int>::size_type i = 0; i < affectedVertexes.size() - 1; i++) {
        for (std::vector<int>::size_type j = i + 1; j < affectedVertexes.size(); j++) {
            Vertex *v1 = affectedVertexes[i];
            Vertex *v2 = affectedVertexes[j];
            Vertex *from, *to = nullptr;
            if (v1->euclideanDistance < v2->euclideanDistance) {
                from = v1;
                to = v2;
            } else {
                from = v2;
                to = v1;
            }
            completeEdges.emplace_back(from, to);
        }
    }
    std::sort(completeEdges.begin(), completeEdges.end(), EdgeWeightComparator());
    // Add all non-intersecting edges to the graph
    for (Edge& completeEdge : completeEdges) {
        std::cout << "Remove consider edge: " << completeEdge << std::endl;
        bool intersects = false;
        for (auto& existingEdge : edges) {
            if (completeEdge.intersects(&existingEdge)) {
                intersects = true;
                std::cout << "  Remove intersects edge: " << existingEdge << std::endl;
                break;
            }
        }
        if (!intersects) {
            auto [it, inserted] = edges.emplace(completeEdge);
            const Edge* edgePtr = &(*it);
            if (inserted) {
                edgePtr->from->edges.emplace_back(edgePtr);
                edgePtr->to->edges.emplace_back(edgePtr);
                std::cout << "  Remove created edge: " << *edgePtr << std::endl;
            } else {
                std::cout << "  Remove skip existing edge: " << *edgePtr << std::endl;
            }
        }
    }
}

void EdgeWeightedGraph::update(Context *ctx) {
    // Delete a vertex and all incoming/outgoing edges if mouse click on
    // Brute force
    // TODO: kd tree
    for (const auto& idVertexPair : vertices) {
        if (idVertexPair.second.inRange(ctx->mouseX, ctx->mouseY)) {
            std::cout << "Removing vertex: " << idVertexPair.second << std::endl;
            removeVertex(idVertexPair.first);
            return;
        }
    }

    // Create a vertex
    auto [idVertexPair, inserted] = vertices.emplace(nextId, Vertex(ctx->mouseX, ctx->mouseY, nextId));
    assert(inserted);
    nextId++;
    Vertex *newVertex = &idVertexPair->second;
    std::cout << "Created vertex: " << *newVertex << std::endl;

    // Create edges for a complete graph
    std::vector<Edge> completeEdges;
    for (auto& vertice : vertices) {
        Vertex *v = &vertice.second;
        if (newVertex->id == v->id) {
            continue;
        }
        Vertex *from, *to = nullptr;
        if (v->euclideanDistance < newVertex->euclideanDistance) {
            from = v;
            to = newVertex;
        } else {
            from = newVertex;
            to = v;
        }
        completeEdges.emplace_back(from, to);
    }
    std::sort(completeEdges.begin(), completeEdges.end(), EdgeWeightComparator());
    // Add all non-intersecting edgers to the graph
    for (Edge& completeEdge : completeEdges) {
        std::cout << "New consider edge: " << completeEdge << std::endl;
        bool intersects = false;
        for (auto& existingEdge : edges) {
            if (completeEdge.intersects(&existingEdge)) {
                intersects = true;
                std::cout << "  New intersects edge: " << existingEdge << std::endl;
                break;
            }
        }
        if (!intersects) {
            auto [it, inserted] = edges.emplace(completeEdge);
            const Edge* edgePtr = &(*it);
            if (inserted) {
                edgePtr->from->edges.emplace_back(edgePtr);
                edgePtr->to->edges.emplace_back(edgePtr);
                std::cout << "  New created edge: " << *edgePtr << std::endl;
            } else {
                std::cout << "  New skip existing edge: " << *edgePtr << std::endl;
            }
        }
    }
}

void EdgeWeightedGraph::render(Context *ctx) {
    // std::cout << "Render Graph" << std::endl;

    // Clear
    SDL_SetRenderDrawColor(ctx->renderer, 255, 255, 255, SDL_ALPHA_OPAQUE);
    SDL_RenderClear(ctx->renderer);

    // Draw the vertices
    for (const auto& idVertexPair : vertices) {
        idVertexPair.second.draw(ctx);
        // std::cout << "  B: " << *v << std::endl;
    }

    // Draw the edges
    for (const auto& idVertexPair : vertices) {
        for (auto edge : idVertexPair.second.edges) {
            if (edge->from->id == idVertexPair.second.id) {
                edge->draw(ctx);
                // std::cout << "    E: " << *edge << std::endl;
            }
        }
    }

    // Present
    SDL_RenderPresent(ctx->renderer);
}
