//
// Created by Sheldon Lobo on 7/6/26.
//

#ifndef MINIMUM_SPANNING_TREE_VERTEX_H
#define MINIMUM_SPANNING_TREE_VERTEX_H

#include <SDL_pixels.h>
#include <vector>

struct Context;
class Edge;

class Vertex {
public:
    const unsigned id;
    int x, y;
    unsigned euclideanDistance;
    std::vector<const Edge*> edges;

    Vertex(unsigned x, unsigned y, unsigned id);
    double distanceTo(Vertex* v) const;
    void removeEdge(const Edge *edge);
    int degree() const;
    void draw(Context* ctx) const;
};
std::ostream& operator<<(std::ostream& strm, const Vertex& v);

class EuclideanDistanceComparator {
public:
    bool operator()(const std::unique_ptr<Vertex>& v1, const std::unique_ptr<Vertex>& v2) const;
};

#endif //MINIMUM_SPANNING_TREE_VERTEX_H