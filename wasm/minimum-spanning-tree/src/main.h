//
// Created by Sheldon Lobo on 7/6/26.
//

#ifndef MINIMUM_SPANNING_TREE_MAIN_H
#define MINIMUM_SPANNING_TREE_MAIN_H

#include <random>
#include <SDL_render.h>
#include <SDL_ttf.h>

extern const unsigned RADIUS;
extern const int DIAMETER;
extern const SDL_Color VERTEX_COLOR;
extern const SDL_Color EDGE_COLOR;
extern const SDL_Color HIGHLIGHT_EDGE_COLOR;

class EdgeWeightedGraph;

struct Context {
    SDL_Renderer *renderer;
    TTF_Font *font;

    bool firstTime;
    unsigned xDimension, yDimension;
    unsigned vertexRadius;
    unsigned sleep;
    bool modified;
    unsigned mouseX, mouseY;

    EdgeWeightedGraph *graph;

    std::mt19937 re;
    std::uniform_int_distribution<> uniformXLocation;
    std::uniform_int_distribution<> uniformYLocation;
};

#endif //MINIMUM_SPANNING_TREE_MAIN_H