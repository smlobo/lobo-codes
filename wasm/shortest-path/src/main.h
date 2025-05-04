//
// Created by Sheldon Lobo on 04/22/25.
//

#ifndef MAIN_H
#define MAIN_H

#include <random>
#include <SDL_render.h>

extern const unsigned RADIUS;
extern const unsigned ARROW;

class EdgeWeightedDigraph;

struct Context {
    SDL_Renderer *renderer;
    bool firstTime;
    unsigned xDimension, yDimension, scale;
    unsigned vertexRadius;
    unsigned sleep;
    bool modified;
    EdgeWeightedDigraph *graph;

    std::mt19937 re;
    std::uniform_int_distribution<> uniformXLocation;
    std::uniform_int_distribution<> uniformYLocation;
};

#endif //MAIN_H
