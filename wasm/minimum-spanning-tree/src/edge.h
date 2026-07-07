//
// Created by Sheldon Lobo on 7/6/26.
//

#ifndef MINIMUM_SPANNING_TREE_EDGE_H
#define MINIMUM_SPANNING_TREE_EDGE_H

#include <iosfwd>
#include <SDL_pixels.h>

struct Context;
class Vertex;

class Edge {
public:
    Vertex* from;
    Vertex* to;
    double weight;
    double angle;
    double angleDegrees;

    // Arrow x & y coordinates
    int xFrom, yFrom, xTo, yTo;

    Edge(Vertex* from, Vertex* to);
    bool intersects(const Edge* edge) const;
    Vertex* other(const Vertex* v) const;
    void draw(Context *ctx, SDL_Color color = SDL_Color{255, 150, 0, SDL_ALPHA_OPAQUE}) const;

    bool operator==(const Edge& other) const;
};
std::ostream& operator<<(std::ostream& strm, const Edge& e);

class EdgeHash {
public:
    std::size_t operator()(const Edge& e) const;
};

class EdgeWeightComparator {
public:
    bool operator()(const Edge& e1, const Edge& e2) const;
};

class EdgeFromComparator {
public:
    bool operator()(const Edge& e1, const Edge& e2) const;
};

#endif //MINIMUM_SPANNING_TREE_EDGE_H