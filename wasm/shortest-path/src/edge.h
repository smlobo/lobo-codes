//
// Created by Sheldon Lobo on 4/25/25.
//

#ifndef EDGE_H
#define EDGE_H

#include <iostream>
#include <SDL2/SDL.h>

struct Context;
class Vertex;

class DirectedEdge {
public:
    Vertex *from;
    Vertex *to;
    double weight;
    double angle;
    double angleDegrees;

    // Arrow x & y coordinates
    int xFrom, yFrom, xTo, yTo;
    // Arrow points
    int xLeft, yLeft, xRight, yRight;

    DirectedEdge(Vertex *from, Vertex *to);
    void draw(Context *ctx, SDL_Color color = SDL_Color{0, 0, 0, SDL_ALPHA_OPAQUE});

    friend std::ostream& operator<<(std::ostream &strm, const DirectedEdge &e);
};

class EdgeWeightComparator {
public:
    bool operator()(const DirectedEdge *e1, const DirectedEdge *e2) const;
};

class EdgeFromComparator {
public:
    bool operator()(const DirectedEdge *e1, const DirectedEdge *e2) const;
};

#endif //EDGE_H
