//
// Created by Sheldon Lobo on 7/25/24.
//

#include <vector>
#include <array>

#include <SDL2/SDL.h>

const int MAX_ITERATIONS = 256;
std::vector<std::array<int, 3>> colors;

void JuliaSetInit() {
    // Generate colors
    for (int col = 0; col < MAX_ITERATIONS; col++) {
        std::array<int, 3> currentColor = {(col >> 5) * 36, (col >> 3 & 7) * 36, (col & 3) * 85};
        colors.push_back(currentColor);
    }
}

void JuliaSet(SDL_Renderer *renderer, int dimension) {
    static float angle = 0;

    // float ca = -0.70176;
    // float cb = -0.3842;
    float ca = cos(angle);
    float cb = sin(angle);
    angle += 0.02;

    int w = 4;
    int h = (w * dimension) / dimension;

    int xmin = -w/2;
    int ymin = -h/2;
    int xmax = xmin + w;
    int ymax = ymin + h;

    float dx = float(xmax - xmin) / dimension;
    float dy = float(ymax - ymin) / dimension;

    float y = ymin;
    for (int j = 0; j < dimension; j++) {
        float x = xmin;
        for (int i = 0; i < dimension; i++) {
            float a = x;
            float b = y;
            int n = 0;
            while (n < MAX_ITERATIONS) {
                float aa = a * a;
                float bb = b * b;
                // Infinite iterations
                if (aa + bb > 64.0) {
                    break;
                }
                float ab2 = 2.0 * a * b;
                a = aa - bb + ca;
                b = ab2 + cb;
                n++;
            }

            // We color each pixel based on how long it takes to get to infinity
            // If we never got there, let's pick the color black (+ offset)
            if (n == MAX_ITERATIONS) {
                SDL_SetRenderDrawColor(renderer, 0, 0, 0, SDL_ALPHA_OPAQUE);
            } else {
                SDL_SetRenderDrawColor(renderer, colors[n][0], colors[n][1], colors[n][2], SDL_ALPHA_OPAQUE);
            }
            SDL_RenderDrawPoint(renderer, i, j);
            x += dx;
        }
        y += dy;
    }
}
