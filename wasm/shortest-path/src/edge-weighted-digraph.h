//
// Created by Sheldon Lobo on 4/25/25.
//

#ifndef EDGE_WEIGHTED_DIGRAPH_H
#define EDGE_WEIGHTED_DIGRAPH_H

#include "vertex.h"

#include <vector>

struct Context;

class EdgeWeightedDigraph {
public:
    std::vector<Vertex> vertices;
    unsigned nEdges;

    explicit EdgeWeightedDigraph(unsigned nVertices, Context& context);
    void render(Context *ctx);
};

#endif //EDGE_WEIGHTED_DIGRAPH_H
