//
// Created by Sheldon Lobo on 6/12/26.
//

#ifndef GRAPH_COLOR_GRAPH_COLORING_H
#define GRAPH_COLOR_GRAPH_COLORING_H

#include "edge-weighted-graph.h"

class SaturationDegree {
public:
    const Vertex* vertex;
    unsigned saturation;

    explicit SaturationDegree(const Vertex* vertex, unsigned saturation = 0);

    bool operator<(const SaturationDegree& other) const;
};

#endif //GRAPH_COLOR_GRAPH_COLORING_H