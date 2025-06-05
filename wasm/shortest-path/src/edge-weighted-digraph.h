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
    std::vector<std::unique_ptr<Vertex>> vertices;
    unsigned nEdges;
    unsigned sourceId, destinationId;

    explicit EdgeWeightedDigraph(unsigned nVertices, Context& context);
    void update(Context *context);
    void render(Context *ctx);
};

#endif //EDGE_WEIGHTED_DIGRAPH_H
