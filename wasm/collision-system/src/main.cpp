//
// Created by Sheldon Lobo on 10/1/24.
//

#include <iostream>

#include <emscripten/emscripten.h>

#include "main.h"
#include "collision-system.h"

constexpr unsigned SLEEP_DEFAULT = 100;
constexpr unsigned SLEEP_DELTA = 5;

void process_input(Context *ctx) {
    SDL_Event event;

    while (SDL_PollEvent(&event)) {
        if ((event.type == SDL_KEYDOWN && event.key.keysym.sym == SDLK_SPACE) ||
            event.type == SDL_MOUSEBUTTONDOWN || event.type == SDL_FINGERDOWN ||
            event.type == SDL_FINGERUP || event.type == SDL_FINGERMOTION) {
            ctx->paused = !ctx->paused;
            std::cout << std::boolalpha << "Paused: " << ctx->paused << "\n";
        } else if (event.type == SDL_KEYDOWN && event.key.keysym.sym == SDLK_LEFT) {
            ctx->sleep += SLEEP_DELTA;
        } else if (event.type == SDL_KEYDOWN && event.key.keysym.sym == SDLK_RIGHT) {
            ctx->sleep -= SLEEP_DELTA;
        } else if (event.type == SDL_KEYDOWN &&
            (event.key.keysym.sym == SDLK_PLUS || event.key.keysym.sym == SDLK_EQUALS ||
                event.key.keysym.sym == SDLK_UP))  {
            std::cout << "Add\n";
            CollisionSystemAdd(ctx);
        } else if (event.type == SDL_KEYDOWN &&
            (event.key.keysym.sym == SDLK_MINUS || event.key.keysym.sym == SDLK_DOWN)) {
            std::cout << "Subtract\n";
            CollisionSystemSubtract(ctx);
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
    ctx->firstTime = false;

    process_input(ctx);

    if (ctx->paused) {
        emscripten_sleep(ctx->sleep);
        return;
    }

    CollisionSystemSimulate(ctx);
    emscripten_sleep(ctx->sleep);
}

extern "C" {

int mainf(int xDim, int yDim) {
    SDL_Window *window;

    Context ctx;
    ctx.firstTime = true;
    ctx.paused = false;
    ctx.xDimension = xDim;
    ctx.yDimension = yDim;
    ctx.scale = std::min(xDim, yDim);
    std::cout << "xDim: " << ctx.xDimension << "; yDim: " << ctx.yDimension << "\n";
    ctx.nBalls = 100;
    ctx.time = 0.0;
    ctx.events = std::priority_queue<Event>();
    ctx.sleep = SLEEP_DEFAULT;

    // Location random generator
    std::random_device rd;
    ctx.re = std::mt19937(rd());
    ctx.uniformLocation = std::uniform_real_distribution<>(0.0,1.0);
    ctx.uniformVelocity = std::uniform_real_distribution<>(-0.005,0.005);

    SDL_Init(SDL_INIT_VIDEO);
    SDL_CreateWindowAndRenderer(ctx.xDimension, ctx.yDimension, 0, &window, &ctx.renderer);
    SDL_SetRenderDrawColor(ctx.renderer, 255, 255, 255, SDL_ALPHA_OPAQUE);
    SDL_RenderClear(ctx.renderer);

    CollisionSystemInit(ctx);

    /**
     * Schedule the main loop handler to get
     * called on each animation frame
     */
    emscripten_set_main_loop_arg(loop_handler, &ctx, 0, true);

    return 0;
}

}