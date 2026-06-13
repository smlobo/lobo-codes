//
// Created by Sheldon Lobo on 6/12/26.
//

#include "graph-coloring.h"

#include <cassert>

SaturationDegree::SaturationDegree(const Vertex* vertex, unsigned saturation) : vertex(vertex), saturation(saturation) {}

bool SaturationDegree::operator<(const SaturationDegree& other) const {
    if (saturation != other.saturation)
        return saturation < other.saturation;
    if (vertex->degree() != other.vertex->degree())
        return vertex->degree() < other.vertex->degree();
    return vertex->id < other.vertex->id;
}

void EdgeWeightedGraph::color() const {
    std::priority_queue<SaturationDegree> dSaturPQ;

    // Reset colors and populate the PQ
    for (const auto& idVertexPair : vertices) {
        assert(idVertexPair.first == idVertexPair.second.id);
        idVertexPair.second.color = -1;
        dSaturPQ.emplace(&idVertexPair.second);
    }

    // Iterate over the PQ
    while (!dSaturPQ.empty()) {
        const Vertex *v = dSaturPQ.top().vertex;
        std::cout << "PQ vertex: " << *v << std::endl;

        // Assign color to the lowest unique color
        std::unordered_set<int> adjVColors;
        for (const auto& e : v->edges) {
            if (e->other(v)->color != -1) {
                adjVColors.insert(e->other(v)->color);
            }
        }
        int color = 0;
        while (!adjVColors.empty() && adjVColors.find(color) != adjVColors.end()) {
            color++;
        }
        std::cout << "  Assigning color: " << color << std::endl;
        v->color = color;

        // Recreate the PQ
        dSaturPQ = std::priority_queue<SaturationDegree>();
        for (const auto& idVertexPair : vertices) {
            v = &idVertexPair.second;
            if (v->color != -1) {
                continue;
            }
            // Compute saturation for this uncolored vertex
            adjVColors.clear();
            for (const auto& e : v->edges) {
                if (e->other(v)->color != -1) {
                    adjVColors.insert(e->other(v)->color);
                }
            }
            dSaturPQ.emplace(v, adjVColors.size());
        }
    }
}
