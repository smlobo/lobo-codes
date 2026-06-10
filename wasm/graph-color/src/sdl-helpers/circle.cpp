//
// Created by Sheldon Lobo on 10/12/24.
//

#include "circle.h"

void DrawCircle(SDL_Renderer * renderer, int centreX, int centreY, int radius, SDL_Color color) {
    SDL_SetRenderDrawColor(renderer, color.r, color.g, color.b, color.a);

    int32_t x = (radius - 1);
    int32_t y = 0;
    int32_t tx = 1;
    int32_t ty = 1;
    int32_t error = (tx - 2*radius);

    while (x >= y)
    {
        //  Each of the following renders an octant of the circle
        SDL_RenderDrawPoint(renderer, centreX + x, centreY - y);
        SDL_RenderDrawPoint(renderer, centreX + x, centreY + y);
        SDL_RenderDrawPoint(renderer, centreX - x, centreY - y);
        SDL_RenderDrawPoint(renderer, centreX - x, centreY + y);
        SDL_RenderDrawPoint(renderer, centreX + y, centreY - x);
        SDL_RenderDrawPoint(renderer, centreX + y, centreY + x);
        SDL_RenderDrawPoint(renderer, centreX - y, centreY - x);
        SDL_RenderDrawPoint(renderer, centreX - y, centreY + x);
        // // Make it 2 pixels
        // SDL_RenderDrawPoint(renderer, centreX + x + 1, centreY - y);
        // SDL_RenderDrawPoint(renderer, centreX + x + 1, centreY + y);
        // SDL_RenderDrawPoint(renderer, centreX - x - 1, centreY - y);
        // SDL_RenderDrawPoint(renderer, centreX - x - 1, centreY + y);
        // SDL_RenderDrawPoint(renderer, centreX + y, centreY - x - 1);
        // SDL_RenderDrawPoint(renderer, centreX + y, centreY + x + 1);
        // SDL_RenderDrawPoint(renderer, centreX - y, centreY - x - 1);
        // SDL_RenderDrawPoint(renderer, centreX - y, centreY + x + 1);

        if (error <= 0)
        {
            ++y;
            error += ty;
            ty += 2;
        }

        if (error > 0)
        {
            --x;
            tx += 2;
            error += (tx - 2*radius);
        }
    }
}

void DrawFilledCircle(SDL_Renderer * renderer, int centreX, int centreY, int radius, SDL_Color color) {
    for (unsigned i = 1; i <= radius; i++) {
        DrawCircle(renderer, centreX, centreY, i, color);
    }
}
