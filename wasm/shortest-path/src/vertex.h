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
    int x, y;
    unsigned euclideanDistance;
    std::set<DirectedEdge *, EdgeWeightComparator> edgesFrom;
    std::set<DirectedEdge *, EdgeWeightComparator> edgesTo;

    Vertex(unsigned x, unsigned y);
    double distanceTo(Vertex &v);
    void draw(Context *ctx, SDL_Color color = SDL_Color{0, 0, 255, SDL_ALPHA_OPAQUE});

    friend std::ostream& operator<<(std::ostream &strm, const Vertex &v);
};

class EuclideanDistanceComparator {
public:
    bool operator()(const Vertex &v1, const Vertex &v2) const;
};

#endif //VERTEX_H
