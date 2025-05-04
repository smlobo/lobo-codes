//
// Created by Sheldon Lobo on 4/25/25.
//

#include "vertex.h"
#include "main.h"
#include "sdl-helpers/circle.h"

#include <iostream>

Vertex::Vertex(unsigned x, unsigned y) : x(x), y(y) {
    euclideanDistance = x*x + y*y;
}

double Vertex::distanceTo(Vertex &v) {
    int xDiff = x - v.x;
    int yDiff = y - v.y;
    return std::sqrt(xDiff * xDiff + yDiff * yDiff);
}

void Vertex::draw(Context *ctx, SDL_Color color) {
    DrawFilledCircle(ctx->renderer, x, y, ctx->vertexRadius, color);
}

std::ostream& operator<<(std::ostream &strm, const Vertex &v) {
    strm << "[" << v.x << "," << v.y << "; From:" << v.edgesFrom.size() << ", To:" << v.edgesTo.size() << "]";
    return strm;
}

bool EuclideanDistanceComparator::operator()(const Vertex &v1, const Vertex &v2) const {
    return v1.euclideanDistance < v2.euclideanDistance;
}
