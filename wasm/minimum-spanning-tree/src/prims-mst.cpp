//
// Created by Sheldon Lobo on 7/10/26.
//

#include "prims-mst.h"

#include <cassert>
#include <ostream>
#include <queue>

#include "edge-weighted-graph.h"

PrimsMinimumSpanningTree::PrimsMinimumSpanningTree(const EdgeWeightedGraph& g) {
    // A MinPQ by edge weight
    std::priority_queue<const Edge*, std::vector<const Edge*>, ReverseEdgeWeightComparator> minEdgePQ;

    // Mark vertices in the MST
    std::unordered_map<unsigned,bool> marked;
    for (const auto& [id, vertex] : g.vertices) {
        marked[id] = false;
    }

    // Add the first vertex edges to the PQ
    const Vertex* first = &g.vertices.begin()->second;
    marked[first->id] = true;
    for (auto& edge : first->edges) {
        minEdgePQ.emplace(edge);
    }

    while (!minEdgePQ.empty()) {
        // Pop the minimum Edge
        const Edge* e = minEdgePQ.top();
        minEdgePQ.pop();

        Vertex* from = e->from;
        Vertex* to = e->to;

        // Ignore if both are in the tree
        if (marked[from->id] && marked[to->id]) {
            continue;
        }

        mstEdges.emplace_back(e);

        Vertex* target = nullptr;
        if (marked[from->id]) {
            assert(!marked[to->id]);
            target = to;
        } else {
            assert(marked[to->id]);
            target = from;
        }

        // Add the target Vertex Edges to the PQ
        marked[target->id] = true;
        for (auto& te : target->edges) {
            if (!marked[te->other(target)->id]) {
                minEdgePQ.emplace(te);
            }
        }
    }
}

std::ostream& operator<<(std::ostream& os, const PrimsMinimumSpanningTree& mst) {
    os << "{";
    for (const Edge* mstEdge : mst.mstEdges) {
        os << *mstEdge << ", ";
    }
    os << "}";
    return os;
}
