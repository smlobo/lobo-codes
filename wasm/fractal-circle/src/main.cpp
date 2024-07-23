#include <iostream>
#include <algorithm>

#include <SDL2/SDL.h>
#include <emscripten/emscripten.h>

// Forward declaration
void FractalCircleInit(SDL_Renderer *, int, int, int, int, bool);

enum InputState
{
    NOTHING_PRESSED = 0,
    LEFT_PRESSED = 1<<0,
    RIGHT_PRESSED = 1<<1
};

struct Context {
    SDL_Renderer *renderer;
    enum InputState activeState;
    bool stateChanged;
    bool firstTime;
    int xDim;
    int yDim;
};

void process_input(struct Context &ctx) {
    SDL_Event event;

    ctx.stateChanged = false;
    while (SDL_PollEvent(&event)) {
        if (event.type == SDL_KEYDOWN) {
            switch (event.key.keysym.sym) {
                case SDLK_LEFT:
                    ctx.activeState = LEFT_PRESSED;
                    ctx.stateChanged = true;
                    break;
                case SDLK_RIGHT:
                    ctx.activeState = RIGHT_PRESSED;
                    ctx.stateChanged = true;
                    break;
                default:
                    ctx.activeState = NOTHING_PRESSED;
                    break;
            }
        } else if (event.type == SDL_MOUSEBUTTONDOWN) {
            switch (event.button.button) {
                case SDL_BUTTON_RIGHT:
                    ctx.activeState = LEFT_PRESSED;
                    ctx.stateChanged = true;
                    break;
                case SDL_BUTTON_LEFT:
                    ctx.activeState = RIGHT_PRESSED;
                    ctx.stateChanged = true;
                    break;
                default:
                    ctx.activeState = NOTHING_PRESSED;
                    break;
            }
        } else if (event.type == SDL_FINGERDOWN) {
            ctx.activeState = RIGHT_PRESSED;
            ctx.stateChanged = true;
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

    if (!ctx->stateChanged) {
        return;
    }

    std::string state = (ctx->activeState == LEFT_PRESSED) ? "Left" : "Right";
    std::cout << "Key state: " << state << "\n";

    SDL_SetRenderDrawColor(ctx->renderer, 255, 255, 255, SDL_ALPHA_OPAQUE);
    SDL_RenderClear(ctx->renderer);

    int circleDiameter = std::min(ctx->xDim, ctx->yDim);
    FractalCircleInit(ctx->renderer, 4, ctx->xDim/2, ctx->yDim/2, circleDiameter, ctx->activeState == RIGHT_PRESSED);

    SDL_RenderPresent(ctx->renderer);
}

extern "C" {

int mainf(int xDim, int yDim) {
    std::cout << "X = " << xDim << ", Y = " << yDim << "\n";

    SDL_Window *window;
    struct Context ctx;
    ctx.firstTime = true;
    ctx.activeState = NOTHING_PRESSED;
    ctx.stateChanged = false;
    ctx.xDim = xDim;
    ctx.yDim = yDim;

    SDL_Init(SDL_INIT_VIDEO);
    SDL_CreateWindowAndRenderer(xDim, yDim, 0, &window, &ctx.renderer);
    SDL_SetRenderDrawColor(ctx.renderer, 255, 255, 255, SDL_ALPHA_OPAQUE);
    SDL_RenderClear(ctx.renderer);

    /**
     * Schedule the main loop handler to get
     * called on each animation frame
     */
    emscripten_set_main_loop_arg(loop_handler, &ctx, -1, 1);

    return 0;
}

}