//
// Created by Sheldon Lobo on 4/25/25.
//

#include "edge-weighted-digraph.h"
#include "main.h"

#include <cassert>
#include <iostream>

EdgeWeightedDigraph::EdgeWeightedDigraph(unsigned nVertices, Context &ctx) : nEdges(0), sourceId(0), destinationId(0),
    nextId(0) {

    // Create random vertices
    for (unsigned i = 0; i < nVertices; i++) {
        vertices.push_back(std::make_unique<Vertex>(ctx.uniformXLocation(ctx.re), ctx.uniformYLocation(ctx.re)));
    }

    // Sort by distance from the origin
    // The shortest is the 'source', the longest is the 'destination' of the graph
    std::sort(vertices.begin(), vertices.end(), EuclideanDistanceComparator());

    // Assign ids to the sorted vector of Vertices
    for (auto &v : vertices) {
        v->setId(nextId++);
    }

    // Set the source and destination Vertex id's
    sourceId = vertices.front()->id;
    destinationId = vertices.back()->id;

    // Draw edges from source to destination
    for (std::vector<int>::size_type i = 0; i < vertices.size() - 1; i++) {
        Vertex *from = vertices[i].get();
        Vertex *to = vertices[i+1].get();
        from->edgesFrom.emplace_back(std::make_shared<DirectedEdge>(from, to));
        to->edgesTo.emplace_back(from->edgesFrom.back());
        std::cout << "Created edge: " << *from->edgesFrom.back() << std::endl;
        nEdges++;
    }

    // Now draw additional edges to vertices at a shorter distance to the existing
    for (std::vector<int>::size_type i = 0; i < vertices.size() - 1; i++) {
        Vertex *vertex = vertices[i].get();
        assert(vertex->edgesFrom.size() == 1);
        double minWeight = (*vertex->edgesFrom.begin())->weight;

        // Iterate over next vertices (NOT prior to avoid cycles)
        for (std::vector<int>::size_type j = i + 2; j < vertices.size(); j++) {
            Vertex *other = vertices[j].get();
            double distance = vertex->distanceTo(other);
            if (distance < minWeight) {
                vertex->edgesFrom.emplace_back(std::make_shared<DirectedEdge>(vertex, other));
                other->edgesTo.emplace_back(vertex->edgesFrom.back());
                std::cout << "Created shorter edge: " << *vertex->edgesFrom.back() << std::endl;
                nEdges++;
            }
        }
    }
}

// Re-assign ids to match the index in the new vertices vector
// Also, update the source & destination id's
void EdgeWeightedDigraph::reassignVertexIds() {
    for (unsigned i = 0; i < vertices.size(); i++) {
        Vertex *vertex = vertices[i].get();
        if (vertex->id == sourceId) {
            sourceId = i;
        } else if (vertex->id == destinationId) {
            destinationId = i;
        }
        vertex->id = i;
    }
}

void EdgeWeightedDigraph::removeVertex(unsigned index) {
    // Remove edges
    Vertex *vertex = vertices[index].get();

    // Remove incoming
    for (unsigned i = 0; i < vertex->edgesTo.size(); i++) {
        DirectedEdge *incomingEdge = vertex->edgesTo[i].get();
        std::cout << "  Removing incoming edge: " << *incomingEdge << std::endl;
        incomingEdge->from->removeOutgoingEdge(incomingEdge);
        nEdges--;
    }
    vertex->edgesTo.clear();
    // Remove outcoming
    for (unsigned i = 0; i < vertex->edgesFrom.size(); i++) {
        DirectedEdge *outgoingEdge = vertex->edgesFrom[i].get();
        std::cout << "  Removing outgoing edge: " << *outgoingEdge << std::endl;
        outgoingEdge->to->removeIncomingEdge(outgoingEdge);
        nEdges--;
    }
    vertex->edgesFrom.clear();

    vertices.erase(vertices.begin() + index);

    reassignVertexIds();
}

void EdgeWeightedDigraph::update(Context *ctx) {
    // Delete a vertex and all incoming/outgoing edges if mouse click on
    // Brute force
    // TODO: kd tree
    for (unsigned i = 0; i < vertices.size(); i++) {
        Vertex *vertex = vertices[i].get();
        if (vertex->inRange(ctx->mouseX, ctx->mouseY)) {
            if (vertex->id == sourceId || vertex->id == destinationId) {
                std::cout << "NOT Removing source/destination vertex: " << *vertex << std::endl;
                return;
            }
            std::cout << "Removing vertex: " << *vertex << std::endl;
            removeVertex(i);
            return;
        }
    }

    // Create a vertex
    // At the end of the vector, but without a conflicting 'origId'
    vertices.push_back(std::make_unique<Vertex>(ctx->mouseX, ctx->mouseY, vertices.size(), nextId++));
    Vertex *newVertex = vertices.back().get();
    std::cout << "Created vertex: " << *newVertex << std::endl;

    // Create an edge from the closest prior vertex to the new vertex
    std::pair<unsigned,double> priorCandidate = {vertices.size(), std::numeric_limits<double>::infinity()};
    std::pair<unsigned,double> subsequentCandidate = {vertices.size(), std::numeric_limits<double>::infinity()};
    for (unsigned i = 0; i < vertices.size(); i++) {
        Vertex *vertex = vertices[i].get();

        // Skip same vertex
        if (newVertex->id == vertex->id) {
            continue;
        }
        double distance = vertex->distanceTo(newVertex);

        // Prior vertices
        if (vertex->euclideanDistance < newVertex->euclideanDistance && distance < priorCandidate.second) {
            priorCandidate = {i, distance};
        }

        // Subsequent vertices
        if (vertex->euclideanDistance > newVertex->euclideanDistance && distance < subsequentCandidate.second) {
            subsequentCandidate = {i, distance};
        }
    }

    // Prior edge
    if (priorCandidate.first < vertices.size()) {
        Vertex *prior = vertices[priorCandidate.first].get();
        // std::cout << "Prior vertex: " << *prior << std::endl;
        prior->edgesFrom.emplace_back(std::make_shared<DirectedEdge>(prior, newVertex));
        newVertex->edgesTo.emplace_back(prior->edgesFrom.back());
        std::cout << "New prior edge: " << *prior->edgesFrom.back() << std::endl;
        nEdges++;
    }

    // Subsequent edge
    if (subsequentCandidate.first < vertices.size()) {
        Vertex *subsequent = vertices[subsequentCandidate.first].get();
        // std::cout << "Subsequent vertex: " << *subsequent << std::endl;
        newVertex->edgesFrom.emplace_back(std::make_shared<DirectedEdge>(newVertex, subsequent));
        subsequent->edgesTo.emplace_back(newVertex->edgesFrom.back());
        std::cout << "New subsequent edge: " << *newVertex->edgesFrom.back() << std::endl;
        nEdges++;
    }
}

void EdgeWeightedDigraph::render(Context *ctx) {
    // std::cout << "Render Graph" << std::endl;

    // Clear
    SDL_SetRenderDrawColor(ctx->renderer, 255, 255, 255, SDL_ALPHA_OPAQUE);
    SDL_RenderClear(ctx->renderer);

    // Draw the vertices
    for (std::vector<int>::size_type i = 0; i < vertices.size(); ++i) {
        Vertex *v = vertices[i].get();
        // Color the source and destination red; the rest are blue
        if (v->id == sourceId) {
            v->draw(ctx, SDL_Color{255, 50, 150, SDL_ALPHA_OPAQUE});
            // std::cout << "  Rsrc: " << *v << std::endl;
        } else if (v->id == destinationId) {
            v->draw(ctx, SDL_Color{255, 150, 50, SDL_ALPHA_OPAQUE});
            // std::cout << "  Rdest: " << *v << std::endl;
        } else {
            v->draw(ctx);
            // std::cout << "  B: " << *v << std::endl;
        }
    }

    // Draw the edges
    for (const auto& v : vertices) {
        for (auto edge : v->edgesFrom) {
            if (ctx->shortestPath->count(edge.get())) {
                edge->draw(ctx, SDL_Color{0, 255, 0, SDL_ALPHA_OPAQUE});
                // std::cout << "    GE: " << *edge << std::endl;
            } else {
                edge->draw(ctx);
                // std::cout << "    E: " << *edge << std::endl;
            }
        }
    }

    // Present
    SDL_RenderPresent(ctx->renderer);
}
