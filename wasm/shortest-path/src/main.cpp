//
// Created by Sheldon Lobo on 4/22/25.
//

#include "main.h"
#include "edge-weighted-digraph.h"
#include "dijkstra-shortest-path.h"

#include <iostream>
#include <emscripten.h>
#include <iomanip>

constexpr unsigned SLEEP_DEFAULT = 100;
constexpr unsigned SLEEP_DELTA = 5;
constexpr unsigned NODES = 10;
constexpr unsigned RADIUS = 15;
constexpr int SEPARATION = RADIUS * 2;
constexpr unsigned ARROW = 10;

void process_input(Context *ctx) {
    SDL_Event event;

    while (SDL_PollEvent(&event)) {
        if (event.type == SDL_MOUSEBUTTONDOWN) {
            ctx->mouseX = event.button.x;
            ctx->mouseY = event.button.y;
            std::cout << std::boolalpha << "Click: " << event.button.x << ", " << event.button.y << "\n";
            ctx->modified = true;
        } else if (event.type == SDL_FINGERDOWN) {
            ctx->mouseX = event.tfinger.x * ctx->xDimension;
            ctx->mouseY = event.tfinger.y * ctx->yDimension;
            std::cout << std::fixed << std::setprecision(4);
            std::cout << std::boolalpha << "Touch: " << event.tfinger.x << "[" << ctx->mouseX << "], " <<
                event.tfinger.y << "[" << ctx->mouseY << "]\n";
            ctx->modified = true;
        }
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

    process_input(ctx);

    if (ctx->modified) {
        // Update the graph with user input - not the first time
        if (!ctx->firstTime) {
            ctx->graph->update(ctx);
        }

        // Calculate the shortest path
        DijkstraShortestPath dsp(ctx->graph);
        ctx->shortestPath = dsp.shortestPath(ctx->graph->destinationId);
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

    ctx->firstTime = false;

    emscripten_sleep(ctx->sleep);
}

extern "C" {

int mainf(int xDim, int yDim) {
    SDL_Window *window;

    Context ctx;
    ctx.firstTime = true;
    ctx.xDimension = xDim;
    ctx.yDimension = yDim;
    std::cout << "xDim: " << ctx.xDimension << "; yDim: " << ctx.yDimension << "\n";
    ctx.sleep = SLEEP_DEFAULT;
    ctx.vertexRadius = RADIUS;
    ctx.modified = true;
    ctx.mouseX = 0;
    ctx.mouseY = 0;

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

    TTF_Init();
    ctx.font = TTF_OpenFont("fonts/SpaceMono-Regular.ttf", 25);
    if (ctx.font == nullptr) {
        std::cout << "TTF_OpenFont error: " << SDL_GetError() << std::endl;
        std::exit(1);
    }

    /**
     * Schedule the main loop handler to get
     * called on each animation frame
     */
    emscripten_set_main_loop_arg(loop_handler, &ctx, 0, true);

    return 0;
}

}