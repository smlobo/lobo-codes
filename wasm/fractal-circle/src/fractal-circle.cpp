//
// Created by Sheldon Lobo on 7/20/24.
//

#include <cmath>
#include <iostream>

#include <SDL2/SDL.h>

static int prevCircles = 0;
static int currCircles = 0;
static int neededCircles = 0;

void DrawCircle(SDL_Renderer * renderer, int centreX, int centreY, int diameter)
{
    if (currCircles == neededCircles) {
        return;
    }

    const int radius = diameter / 2;

    int32_t x = (radius - 1);
    int32_t y = 0;
    int32_t tx = 1;
    int32_t ty = 1;
    int32_t error = (tx - diameter);

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
        // Make it 2 pixels
        SDL_RenderDrawPoint(renderer, centreX + x + 1, centreY - y);
        SDL_RenderDrawPoint(renderer, centreX + x + 1, centreY + y);
        SDL_RenderDrawPoint(renderer, centreX - x - 1, centreY - y);
        SDL_RenderDrawPoint(renderer, centreX - x - 1, centreY + y);
        SDL_RenderDrawPoint(renderer, centreX + y, centreY - x - 1);
        SDL_RenderDrawPoint(renderer, centreX + y, centreY + x + 1);
        SDL_RenderDrawPoint(renderer, centreX - y, centreY - x - 1);
        SDL_RenderDrawPoint(renderer, centreX - y, centreY + x + 1);

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
            error += (tx - diameter);
        }
    }

    currCircles++;
}

void FractalCircle(SDL_Renderer *renderer, int depth, int centerX, int centerY, int diameter) {
    // Base case
    if (depth == 0)
        return;

    // SDL_SetRenderDrawColor(renderer, randomInt(0, 256), randomInt(0, 256), randomInt(0, 256), SDL_ALPHA_OPAQUE);
    SDL_SetRenderDrawColor(renderer, 255, 0, 0, SDL_ALPHA_OPAQUE);

    // Draw circle
    DrawCircle(renderer, centerX, centerY, diameter);

    // Draw 7 circles
    int innerD = diameter / 3;
    // cos 30
    int xOffset = static_cast<int>(static_cast<float>(innerD) * sqrt(3.0) / 2.0);
    // sin 30
    int yOffset = innerD / 2;

    // Center
    FractalCircle(renderer, depth-1, centerX, centerY, innerD);
    // Above
    FractalCircle(renderer, depth-1, centerX, centerY-innerD, innerD);
    // Below
    FractalCircle(renderer, depth-1, centerX, centerY+innerD, innerD);
    // Above left
    FractalCircle(renderer, depth-1, centerX-xOffset, centerY-yOffset, innerD);
    // Above right
    FractalCircle(renderer, depth-1, centerX+xOffset, centerY-yOffset, innerD);
    // Below left
    FractalCircle(renderer, depth-1, centerX-xOffset, centerY+yOffset, innerD);
    // Below right
    FractalCircle(renderer, depth-1, centerX+xOffset, centerY+yOffset, innerD);
}

void FractalCircleInit(SDL_Renderer *renderer, int depth, int centerX, int centerY, int diameter, bool forward) {
    if (forward) {
        neededCircles = prevCircles + 1;
    } else {
        if (prevCircles == 0) {
            return;
        }
        neededCircles = prevCircles - 1;
    }
    currCircles = 0;
    FractalCircle(renderer, depth, centerX, centerY, diameter);

    // Update circles on the screen if we drew what we wanted to
    if (currCircles == neededCircles) {
        prevCircles = neededCircles;
    }

    std::cout << "Circles drawn: " << prevCircles << "\n";
}