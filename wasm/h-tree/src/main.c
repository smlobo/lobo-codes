#include <SDL2/SDL.h>
#include <emscripten.h>

void HTree(SDL_Renderer *renderer, int depth, int size, int centerX, int centerY);

void mainf(int dimension) {
    SDL_Window *window;
    SDL_Renderer *renderer;

    SDL_Init(SDL_INIT_VIDEO);

    SDL_CreateWindowAndRenderer(dimension, dimension, 0, &window, &renderer);
    SDL_SetRenderDrawColor(renderer, 255, 255, 255, 255);
    SDL_RenderClear(renderer);

    HTree(renderer, 5, dimension/2, dimension/2, dimension/2);

    SDL_RenderPresent(renderer);
}
