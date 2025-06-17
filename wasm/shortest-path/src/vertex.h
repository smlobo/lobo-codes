//
// Created by Sheldon Lobo on 4/25/25.
//

#ifndef VERTEX_H
#define VERTEX_H

#include "edge.h"

#include <set>
#include <SDL2/SDL.h>

class EdgeWeightComparator;
struct Context;

class Vertex {
public:
    unsigned id, origId;
    int x, y;
    unsigned euclideanDistance;
    std::vector<std::shared_ptr<DirectedEdge>> edgesFrom;
    std::vector<std::shared_ptr<DirectedEdge>> edgesTo;

    Vertex(unsigned x, unsigned y, unsigned id = 0, unsigned oridId = 0);
    void setId(unsigned id);
    double distanceTo(Vertex *v) const;
    bool inRange(int x, int y) const;
    void removeIncomingEdge(const DirectedEdge *edge);
    void removeOutgoingEdge(const DirectedEdge *edge);
    void draw(Context *ctx, SDL_Color color = SDL_Color{150, 150, 255, SDL_ALPHA_OPAQUE});

    friend std::ostream& operator<<(std::ostream &strm, const Vertex &v);
};

class EuclideanDistanceComparator {
public:
    bool operator()(const std::unique_ptr<Vertex>& v1, const std::unique_ptr<Vertex>& v2) const;
};

#endif //VERTEX_H
