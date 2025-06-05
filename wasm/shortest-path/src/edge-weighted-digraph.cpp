//
// Created by Sheldon Lobo on 4/25/25.
//

#include "edge-weighted-digraph.h"
#include "main.h"

#include <cassert>
#include <iostream>

EdgeWeightedDigraph::EdgeWeightedDigraph(unsigned nVertices, Context &ctx) : nEdges(0), sourceId(0), destinationId(0) {

    // Create random vertices
    for (unsigned i = 0; i < nVertices; i++) {
        vertices.push_back(std::make_unique<Vertex>(ctx.uniformXLocation(ctx.re), ctx.uniformYLocation(ctx.re)));
    }

    // Sort by distance from the origin
    // The shortest is the 'source', the longest is the 'destination' of the graph
    std::sort(vertices.begin(), vertices.end(), EuclideanDistanceComparator());

    // Assign ids to the sorted vector of Vertices
    for (unsigned i = 0; i < vertices.size(); i++) {
        vertices[i]->id = i;
    }

    // Set the source and destination Vertex id's
    sourceId = vertices.front()->id;
    destinationId = vertices.back()->id;

    // Draw edges from source to destination
    for (std::vector<int>::size_type i = 0; i < vertices.size() - 1; i++) {
        Vertex *from = vertices[i].get();
        Vertex *to = vertices[i+1].get();
        DirectedEdge *edge = new DirectedEdge(from, to);
        from->edgesFrom.insert(edge);
        to->edgesTo.insert(edge);
        // std::cout << "Created edge: " << *from << " -> " << to << std::endl;
        nEdges++;
    }

    // Now draw additional edges to vertices at a shorter distance to the existing
    for (std::vector<int>::size_type i = 0; i < vertices.size() - 1; i++) {
        Vertex *vertex = vertices[i].get();
        assert(vertex->edgesFrom.size() == 1);
        DirectedEdge *edge = *vertex->edgesFrom.begin();
        double minWeight = edge->weight;

        // Iterate over next vertices (NOT prior to avoid cycles)
        for (std::vector<int>::size_type j = i + 2; j < vertices.size(); j++) {
            Vertex *other = vertices[j].get();
            double distance = vertex->distanceTo(other);
            if (distance < minWeight) {
                auto *newEdge = new DirectedEdge(vertex, other);
                vertex->edgesFrom.insert(newEdge);
                other->edgesTo.insert(newEdge);
                // std::cout << "Created shorter edge: " << i << " " << *vertex << " -> " << j << " " << other << std::endl;
                nEdges++;
            }
        }
    }
}

void EdgeWeightedDigraph::update(Context *ctx) {
    // Create a vertex
    // TODO: Delete if vertex at same coords
    vertices.push_back(std::make_unique<Vertex>(ctx->mouseX, ctx->mouseY, vertices.size()));
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
        DirectedEdge *priorEdge = new DirectedEdge(prior, newVertex);
        prior->edgesFrom.insert(priorEdge);
        newVertex->edgesTo.insert(priorEdge);
        std::cout << "New prior edge: " << *priorEdge << std::endl;
        nEdges++;
    }

    // Subsequent edge
    if (subsequentCandidate.first < vertices.size()) {
        Vertex *subsequent = vertices[subsequentCandidate.first].get();
        // std::cout << "Subsequent vertex: " << *subsequent << std::endl;
        DirectedEdge *subsequentEdge = new DirectedEdge(newVertex, subsequent);
        newVertex->edgesFrom.insert(subsequentEdge);
        subsequent->edgesTo.insert(subsequentEdge);
        std::cout << "New subsequent edge: " << *subsequentEdge << std::endl;
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
            if (ctx->shortestPath->count(edge)) {
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
