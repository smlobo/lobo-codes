//
// Created by Sheldon Lobo on 5/14/25.
//

#include "dijkstra-shortest-path.h"

#include <iomanip>

DijkstraShortestPath::DijkstraShortestPath(EdgeWeightedDigraph *g) : g(g) {
    std::priority_queue<RelaxEdge, std::vector<RelaxEdge>, RelaxEdgeComparator> pq;

    edgeTo = std::vector<DirectedEdge*>(g->vertices.size(), nullptr);
    distTo = std::vector<double>(g->vertices.size(), std::numeric_limits<double>::infinity());

    distTo[g->sourceId] = 0.0;

    pq.emplace(g->sourceId, 0.0);
    while (!pq.empty()) {
        RelaxEdge rEdge = pq.top();
        pq.pop();
        unsigned toId = rEdge.v;
        for (DirectedEdge *e : g->vertices[toId].get()->edgesFrom)
            relax(e, pq);
    }

    std::cout << "EdgeTo: ";
    for (unsigned i = 0; i < edgeTo.size(); i++) {
        std::cout << "[" << edgeTo[i]->from->id << " -> " << i << "]; " ;
    }
    std::cout << "\n";
    // std::cout << "DistTo:\n";
    // std::cout << std::fixed << std::setprecision(4);
    // for (unsigned i = 0; i < distTo.size(); i++) {
    //     std::cout << "    [" << i << "] " << distTo[i] << "\n" ;
    // }
}

std::set<DirectedEdge*, EdgeFromComparator> *DijkstraShortestPath::shortestPath(unsigned d) const {
    auto *shortestPath = new std::set<DirectedEdge*, EdgeFromComparator>();
    std::cout << "Shortest Path: ";

    unsigned v = d;
    while (v != 0) {
        std::cout << v << "/" << g->vertices[v]->origId << ", ";
        // Impossible
        if (edgeTo[v] == nullptr) {
            std::cout << "Impossible edge to: " << v << std::endl;
            break;
        }
        shortestPath->insert(edgeTo[v]);
        v = edgeTo[v]->from->id;
    }
    std::cout << v << std::endl;

    return shortestPath;
}

void DijkstraShortestPath::relax(DirectedEdge *e, std::priority_queue<RelaxEdge, std::vector<RelaxEdge>,
    RelaxEdgeComparator> &pq) {

    unsigned v = e->from->id;
    unsigned w = e->to->id;

    if (distTo[w] > distTo[v] + e->weight) {

        distTo[w] = distTo[v] + e->weight;
        edgeTo[w] = e;

        pq.emplace(e->to->id, distTo[w]);
    }

}

RelaxEdge::RelaxEdge(unsigned v, double w) : v(v), weight(w) {}

bool RelaxEdgeComparator::operator()(const RelaxEdge &e1, const RelaxEdge &e2) const {
    return e1.weight > e2.weight;
}
