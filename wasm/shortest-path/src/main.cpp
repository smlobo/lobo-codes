//
// Created by Sheldon Lobo on 4/22/25.
//

#include "main.h"
#include "edge-weighted-digraph.h"
#include "dijkstra-shortest-path.h"

#include <iostream>
#include <emscripten/emscripten.h>
#include <SDL.h>

constexpr unsigned SLEEP_DEFAULT = 100;
constexpr unsigned SLEEP_DELTA = 5;
constexpr unsigned NODES = 10;
constexpr unsigned RADIUS = 10;
constexpr unsigned ARROW = 10;

void process_input(Context *ctx) {
    SDL_Event event;

    while (SDL_PollEvent(&event)) {
    }
}

void loop_handler(void *arg)
{
    auto *ctx = static_cast<Context *>(arg);

    if (ctx->firstTime) {
        SDL_SetRenderDrawColor(ctx->renderer, 255, 255, 255, SDL_ALPHA_OPAQUE);
        SDL_RenderClear(ctx->renderer);
        SDL_RenderPresent(ctx->renderer);
    }
    ctx->firstTime = false;

    process_input(ctx);

    if (ctx->modified) {
        // Calculate the shortest path
        DijkstraShortestPath dsp(ctx->graph);
        ctx->shortestPath = dsp.shortestPath(ctx->graph->vertices.size() - 1);
        // Print to console
        std::cout << "Shortest Path:\n";
        for (DirectedEdge *e : *ctx->shortestPath) {
            std::cout << "  " << *e << "\n";
        }
        ctx->graph->render(ctx);
        delete ctx->shortestPath;
        ctx->shortestPath = nullptr;
        ctx->modified = false;
    }

    emscripten_sleep(ctx->sleep);
}

extern "C" {

int mainf(int xDim, int yDim) {
    SDL_Window *window;

    Context ctx;
    ctx.firstTime = true;
    ctx.xDimension = xDim;
    ctx.yDimension = yDim;
    ctx.scale = std::min(xDim, yDim);
    std::cout << "xDim: " << ctx.xDimension << "; yDim: " << ctx.yDimension << "\n";
    ctx.sleep = SLEEP_DEFAULT;
    ctx.vertexRadius = RADIUS;
    ctx.modified = true;

    // Location random generator
    std::random_device rd;
    ctx.re = std::mt19937(rd());
    ctx.uniformXLocation = std::uniform_int_distribution<>(0 + RADIUS,ctx.xDimension - RADIUS);
    ctx.uniformYLocation = std::uniform_int_distribution<>(0 + RADIUS,ctx.yDimension - RADIUS);

    SDL_Init(SDL_INIT_VIDEO);
    SDL_CreateWindowAndRenderer(ctx.xDimension, ctx.yDimension, 0, &window, &ctx.renderer);
    SDL_SetRenderDrawColor(ctx.renderer, 255, 255, 255, SDL_ALPHA_OPAQUE);
    SDL_RenderClear(ctx.renderer);

    EdgeWeightedDigraph edgeWeightedDigraph(NODES, ctx);
    ctx.graph = &edgeWeightedDigraph;
    ctx.shortestPath = nullptr;
    // edgeWeightedDigraph.render(&ctx);

    /**
     * Schedule the main loop handler to get
     * called on each animation frame
     */
    emscripten_set_main_loop_arg(loop_handler, &ctx, 0, true);

    return 0;
}

}