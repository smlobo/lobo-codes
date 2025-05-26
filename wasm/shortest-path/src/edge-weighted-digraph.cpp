//
// Created by Sheldon Lobo on 4/25/25.
//

#include "edge-weighted-digraph.h"
#include "main.h"

#include <cassert>
#include <iostream>

EdgeWeightedDigraph::EdgeWeightedDigraph(unsigned nVertices, Context &ctx) {

    nEdges = 0;

    // Create random vertices
    for (unsigned i = 0; i < nVertices; i++) {
        vertices.emplace_back(ctx.uniformXLocation(ctx.re), ctx.uniformYLocation(ctx.re));
    }

    // Sort by distance from the origin
    // The shortest is the 'source', the longest is the 'destination' of the graph
    std::sort(vertices.begin(), vertices.end(), EuclideanDistanceComparator());

    // Assign ids to the sorted verctor of Vertices
    for (unsigned i = 0; i < vertices.size(); i++) {
        vertices[i].id = i;
    }

    // Draw edges from source to destination
    for (std::vector<int>::size_type i = 0; i < vertices.size() - 1; i++) {
        Vertex &from = vertices[i];
        Vertex &to = vertices[i+1];
        DirectedEdge *edge = new DirectedEdge(&from, &to);
        from.edgesFrom.insert(edge);
        to.edgesTo.insert(edge);
        std::cout << "Created edge: " << i << " " << from << " -> " << (i + 1) << " " << to << std::endl;
        nEdges++;
    }

    // Now draw additional edges to vertices at a shorter distance to the existing
    for (std::vector<int>::size_type i = 0; i < vertices.size() - 1; i++) {
        Vertex &vertex = vertices[i];
        assert(vertex.edgesFrom.size() == 1);
        DirectedEdge *edge = *vertex.edgesFrom.begin();
        double minWeight = edge->weight;

        // Iterate over next vertices (NOT prior to avoid cycles)
        for (std::vector<int>::size_type j = i + 2; j < vertices.size(); j++) {
            Vertex &other = vertices[j];
            double distance = vertex.distanceTo(other);
            if (distance < minWeight) {
                auto *newEdge = new DirectedEdge(&vertex, &other);
                vertex.edgesFrom.insert(newEdge);
                other.edgesTo.insert(newEdge);
                std::cout << "Created shorter edge: " << i << " " << vertex << " -> " << j << " " << other << std::endl;
                nEdges++;
            }
        }
    }
}

void EdgeWeightedDigraph::render(Context *ctx) {
    std::cout << "Render Graph" << std::endl;

    // Clear
    SDL_SetRenderDrawColor(ctx->renderer, 255, 255, 255, SDL_ALPHA_OPAQUE);
    SDL_RenderClear(ctx->renderer);

    // Draw the vertices
    for (std::vector<int>::size_type i = 0; i < vertices.size(); ++i) {
        Vertex &v = vertices[i];
        // Color the source and destination red; the rest are blue
        if (i == 0 || i == vertices.size() - 1) {
            v.draw(ctx, SDL_Color{255, 0, 0, SDL_ALPHA_OPAQUE});
            std::cout << "  R: {" << i << "} " << v << std::endl;
        } else {
            v.draw(ctx);
            std::cout << "  B: {" << i << "} " << v << std::endl;
        }
    }

    // Draw the edges
    for (auto v : vertices) {
        for (auto edge : v.edgesFrom) {
            if (ctx->shortestPath->count(edge)) {
                edge->draw(ctx, SDL_Color{0, 255, 0, SDL_ALPHA_OPAQUE});
                std::cout << "    GE: " << *edge << std::endl;
            } else {
                edge->draw(ctx);
                std::cout << "    E: " << *edge << std::endl;
            }
        }
    }

    // Present
    SDL_RenderPresent(ctx->renderer);
}
