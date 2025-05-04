//
// Created by Sheldon Lobo on 4/25/25.
//

#ifndef EDGE_WEIGHTED_DIGRAPH_H
#define EDGE_WEIGHTED_DIGRAPH_H

#include "vertex.h"

#include <vector>

struct Context;
// class Vertex;

class EdgeWeightedDigraph {
private:
    std::vector<Vertex> vertices;
    unsigned nEdges;
public:
    explicit EdgeWeightedDigraph(unsigned nVertices, Context& context);
    void render(Context *ctx);
};

#endif //EDGE_WEIGHTED_DIGRAPH_H
