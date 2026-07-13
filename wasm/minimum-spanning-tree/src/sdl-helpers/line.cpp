//
// Created by Sheldon Lobo on 7/12/26.
//

#include "line.h"

#include <algorithm>
#include <cmath>

void DrawLine(SDL_Renderer* renderer, int xFrom, int yFrom, int xTo, int yTo, int width, SDL_Color color) {
    SDL_SetRenderDrawColor(renderer, color.r, color.g, color.b, color.a);

    const int lineWidth = std::max(1, width);
    if (lineWidth == 1) {
        SDL_RenderDrawLine(renderer, xFrom, yFrom, xTo, yTo);
        return;
    }

    const double xDiff = xTo - xFrom;
    const double yDiff = yTo - yFrom;
    const double length = std::sqrt(xDiff * xDiff + yDiff * yDiff);
    if (length == 0) {
        SDL_RenderDrawPoint(renderer, xFrom, yFrom);
        return;
    }

    const double normalX = -yDiff / length;
    const double normalY = xDiff / length;
    const double centerOffset = (lineWidth - 1) / 2.0;

    for (int i = 0; i < lineWidth; ++i) {
        const double offset = i - centerOffset;
        const int xOffset = static_cast<int>(std::round(normalX * offset));
        const int yOffset = static_cast<int>(std::round(normalY * offset));
        SDL_RenderDrawLine(renderer, xFrom + xOffset, yFrom + yOffset, xTo + xOffset, yTo + yOffset);
    }
}
