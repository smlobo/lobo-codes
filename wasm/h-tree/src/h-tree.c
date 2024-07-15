//
// Created by Sheldon Lobo on 7/7/24.
//

#include <stdbool.h>
#include <SDL2/SDL.h>

void DrawLine(SDL_Renderer *renderer, int x1, int y1, int x2, int y2, int width) {
    // Slope
    int xDiff = abs(x2 - x1);
    int yDiff = abs(y2 - y1);
    bool horizontal = xDiff > yDiff;
    for (int i = -width/2; i < (width+1)/2; i++) {
        horizontal ? SDL_RenderDrawLine(renderer, x1, y1 + i, x2, y2 + i) :
        SDL_RenderDrawLine(renderer, x1 + i, y1, x2 + i, y2);
    }
}

void HTree(SDL_Renderer *renderer, int depth, int size, int centerX, int centerY) {
    // Base case
    if (depth == 0)
        return;

    SDL_SetRenderDrawColor(renderer, 255, 0, 0, SDL_ALPHA_OPAQUE);

    // Draw H to size and centered

    // Horizontal
    DrawLine(renderer, centerX - size/2, centerY, centerX + size/2, centerY, 2);
    // Vertical left
    DrawLine(renderer, centerX - size/2, centerY - size/2, centerX - size/2, centerY + size/2, 2);
    // Horizontal
    DrawLine(renderer, centerX + size/2, centerY - size/2, centerX + size/2, centerY + size/2, 2);

    // Draw 4 H's at the 4 corners
    HTree(renderer, depth-1, size/2, centerX - size/2, centerY - size/2);
    HTree(renderer, depth-1, size/2, centerX + size/2, centerY - size/2);
    HTree(renderer, depth-1, size/2, centerX - size/2, centerY + size/2);
    HTree(renderer, depth-1, size/2, centerX + size/2, centerY + size/2);
}
