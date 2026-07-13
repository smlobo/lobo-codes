//
// Created by Sheldon Lobo on 7/10/26.
//

#ifndef MINIMUM_SPANNING_TREE_PRIMS_MST_H
#define MINIMUM_SPANNING_TREE_PRIMS_MST_H

#include <vector>

#include "edge.h"

class EdgeWeightedGraph;

class PrimsMinimumSpanningTree {
public:
    std::vector<const Edge*> mstEdges;

    PrimsMinimumSpanningTree(const EdgeWeightedGraph& g);
};
std::ostream& operator<<(std::ostream& os, const PrimsMinimumSpanningTree& g);

#endif //MINIMUM_SPANNING_TREE_PRIMS_MST_H