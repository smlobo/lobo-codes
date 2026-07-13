//
// Created by Sheldon Lobo on 7/6/26.
//

#ifndef MINIMUM_SPANNING_TREE_EDGE_WEIGHTED_GRAPH_H
#define MINIMUM_SPANNING_TREE_EDGE_WEIGHTED_GRAPH_H

#include <unordered_map>
#include <unordered_set>

#include "cell.h"
#include "edge.h"
#include "vertex.h"

class EdgeHash;
struct Context;

class EdgeWeightedGraph {
public:
    std::unordered_map<unsigned, Vertex> vertices;
    std::unordered_map<Cell, unsigned, CellHash> cellVertexIdMap;
    std::unordered_set<Edge, EdgeHash> edges;
    unsigned nextId;

    EdgeWeightedGraph(unsigned nVertices, Context& context);
    void update(Context* context);
    void removeVertex(unsigned id);

    void computeMinimumSpanningTree();

    Cell cellVertex(int x, int y) const;
    Cell cellVertexOrNeighbor(int x, int y) const;
    Cell cellVertexOrCornerNeighbor(int x, int y) const;

    void draw(Context* ctx) const;
};

#endif //MINIMUM_SPANNING_TREE_EDGE_WEIGHTED_GRAPH_H