//
// Created by Sheldon Lobo on 7/6/26.
//

#ifndef MINIMUM_SPANNING_TREE_EDGE_WEIGHTED_GRAPH_H
#define MINIMUM_SPANNING_TREE_EDGE_WEIGHTED_GRAPH_H

#include <map>
#include <unordered_set>

#include "edge.h"
#include "vertex.h"

class EdgeHash;
struct Context;

class EdgeWeightedGraph {
public:
    std::map<unsigned, Vertex> vertices;
    std::unordered_set<Edge, EdgeHash> edges;
    unsigned nextId;

    EdgeWeightedGraph(unsigned nVertices, Context& context);
    void update(Context* context);
    void removeVertex(unsigned id);
    void render(Context* ctx);
};

#endif //MINIMUM_SPANNING_TREE_EDGE_WEIGHTED_GRAPH_H