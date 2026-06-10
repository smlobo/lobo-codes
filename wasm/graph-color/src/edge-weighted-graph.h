//
// Created by Sheldon Lobo on 5/29/26.
//

#ifndef EDGE_WEIGHTED_DIGRAPH_H
#define EDGE_WEIGHTED_DIGRAPH_H

#include <map>
#include <unordered_set>

#include "vertex.h"

struct Context;

class EdgeWeightedGraph {
public:
    std::map<unsigned, Vertex> vertices;
    std::unordered_set<Edge, EdgeHash> edges;
    unsigned nextId;

    explicit EdgeWeightedGraph(unsigned nVertices, Context& context);
    void update(Context *context);
    void removeVertex(unsigned id);
    void render(Context *ctx);
};

#endif //EDGE_WEIGHTED_DIGRAPH_H
