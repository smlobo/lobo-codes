//
// Created by Sheldon Lobo on 5/14/25.
//

#ifndef DIJKSTRA_SHORTEST_PATH_H
#define DIJKSTRA_SHORTEST_PATH_H

#include "edge-weighted-digraph.h"
#include "vertex.h"

class RelaxEdgeComparator;
class RelaxEdge;

class DijkstraShortestPath {
private:
    void relax(DirectedEdge *e, std::priority_queue<RelaxEdge, std::vector<RelaxEdge>, RelaxEdgeComparator> &pq);
public:
    EdgeWeightedDigraph *g;
    std::vector<DirectedEdge*> edgeTo;
    std::vector<double> distTo;

    explicit DijkstraShortestPath(EdgeWeightedDigraph *g);
    std::set<DirectedEdge*, EdgeFromComparator> *shortestPath(unsigned d) const;
};

class RelaxEdge {
public:
    unsigned v;
    double weight;

    RelaxEdge(unsigned v, double w);
};

class RelaxEdgeComparator {
public:
    bool operator()(const RelaxEdge &e1, const RelaxEdge &e2) const;
};

#endif //DIJKSTRA_SHORTEST_PATH_H
