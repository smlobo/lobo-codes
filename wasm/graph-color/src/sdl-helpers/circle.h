//
// Created by Sheldon Lobo on 10/12/24.
//

#ifndef CIRCLE_H
#define CIRCLE_H

#include <SDL_render.h>

void DrawCircle(SDL_Renderer *renderer, int centreX, int centreY, int radius, SDL_Color color);
void DrawFilledCircle(SDL_Renderer *renderer, int centreX, int centreY, int radius, SDL_Color color);

#endif //CIRCLE_H
