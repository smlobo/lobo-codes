//
// Created by Sheldon Lobo on 7/25/24.
//

#include <iostream>
#include <thread>

#include <SDL2/SDL.h>
#include <emscripten/emscripten.h>

const int MAX_DIMENSION = 1000;

// Forward declaration
void JuliaSetInit();
void JuliaSet(SDL_Renderer *, int);

struct Context {
    SDL_Renderer *renderer;
    bool paused;
    bool firstTime;
    int dimension;
};

void process_input(struct Context &ctx) {
    SDL_Event event;

    while (SDL_PollEvent(&event)) {
    // while (SDL_WaitEventTimeout(&event, 10)) {
        if ((event.type == SDL_KEYDOWN && event.key.keysym.sym == SDLK_SPACE) ||
            event.type == SDL_MOUSEBUTTONDOWN || event.type == SDL_FINGERDOWN ||
            event.type == SDL_FINGERUP || event.type == SDL_FINGERMOTION) {
            ctx.paused = !ctx.paused;
            std::cout << std::boolalpha << "Paused: " << ctx.paused << "\n";
        }
    }
}

void loop_handler(void *arg)
{
    auto *ctx = static_cast<struct Context *>(arg);

    if (ctx->firstTime) {
        SDL_SetRenderDrawColor(ctx->renderer, 255, 255, 255, SDL_ALPHA_OPAQUE);
        SDL_RenderClear(ctx->renderer);
        SDL_RenderPresent(ctx->renderer);
    }
    ctx->firstTime = false;

    process_input(*ctx);

    if (ctx->paused) {
        return;
    }

    SDL_SetRenderDrawColor(ctx->renderer, 255, 255, 255, SDL_ALPHA_OPAQUE);
    SDL_RenderClear(ctx->renderer);

    JuliaSet(ctx->renderer, ctx->dimension);

    SDL_RenderPresent(ctx->renderer);

    emscripten_sleep(100);
}

extern "C" {

int mainf(int xDim, int yDim) {
    SDL_Window *window;

    Context ctx;
    ctx.firstTime = true;
    ctx.paused = false;
    ctx.dimension = std::min(xDim, yDim);
    ctx.dimension = std::min(ctx.dimension, MAX_DIMENSION);
    std::cout << "xDim: " << xDim << "; yDim: " << yDim << "; MAX: " << MAX_DIMENSION <<
        "; Set to: " << ctx.dimension << "\n";

    JuliaSetInit();

    SDL_Init(SDL_INIT_VIDEO);
    SDL_CreateWindowAndRenderer(ctx.dimension, ctx.dimension, 0, &window, &ctx.renderer);
    SDL_SetRenderDrawColor(ctx.renderer, 255, 255, 255, SDL_ALPHA_OPAQUE);
    SDL_RenderClear(ctx.renderer);

    /**
     * Schedule the main loop handler to get
     * called on each animation frame
     */
    emscripten_set_main_loop_arg(loop_handler, &ctx, 0, true);

    return 0;
}

}