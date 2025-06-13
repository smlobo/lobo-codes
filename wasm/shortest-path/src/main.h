//
// Created by Sheldon Lobo on 4/22/25.
//

#ifndef MAIN_H
#define MAIN_H

#include <random>
#include <SDL_render.h>
#include <SDL_ttf.h>
#include <set>

extern const unsigned RADIUS;
extern const int SEPARATION;
extern const unsigned ARROW;

class EdgeWeightedDigraph;
class DirectedEdge;
class EdgeFromComparator;

struct Context {
    SDL_Renderer *renderer;
    TTF_Font *font;

    bool firstTime;
    unsigned xDimension, yDimension;
    unsigned vertexRadius;
    unsigned sleep;
    bool modified;
    unsigned mouseX, mouseY;

    EdgeWeightedDigraph *graph;
    std::set<DirectedEdge*, EdgeFromComparator> *shortestPath;

    std::mt19937 re;
    std::uniform_int_distribution<> uniformXLocation;
    std::uniform_int_distribution<> uniformYLocation;
};

#endif //MAIN_H
