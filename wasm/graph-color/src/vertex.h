//
// Created by Sheldon Lobo on 5/29/26.
//

#ifndef VERTEX_H
#define VERTEX_H

#include "edge.h"

#include <set>
#include <SDL2/SDL.h>

struct Context;
class Edge;

class Vertex {
public:
    const unsigned id;
    int x, y;
    unsigned euclideanDistance;
    std::vector<const Edge*> edges;

    Vertex(unsigned x, unsigned y, unsigned id);
    double distanceTo(Vertex *v) const;
    bool inRange(int x, int y) const;
    bool tooClose(int x, int y) const;
    void removeEdge(const Edge *edge);
    void draw(Context *ctx, SDL_Color color = SDL_Color{150, 255, 150, SDL_ALPHA_OPAQUE}) const;

    bool operator<(const Vertex &v) const;
    friend std::ostream& operator<<(std::ostream &strm, const Vertex &v);
};

class EuclideanDistanceComparator {
public:
    bool operator()(const std::unique_ptr<Vertex>& v1, const std::unique_ptr<Vertex>& v2) const;
};

#endif //VERTEX_H
