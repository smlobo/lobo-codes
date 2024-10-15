//
// Created by Sheldon Lobo on 10/6/24.
//

#ifndef MAIN_H
#define MAIN_H

#include <queue>
#include <random>

#include <SDL2/SDL.h>

#include "event.h"

class Ball;

struct Context {
    SDL_Renderer *renderer;
    bool paused;
    bool firstTime;
    unsigned xDimension, yDimension, scale;
    unsigned nBalls;
    std::deque<Ball*> balls;
    double time;
    std::priority_queue<Event> events;
    unsigned sleep;

    std::mt19937 re;
    std::uniform_real_distribution<> uniformLocation;
    std::uniform_real_distribution<> uniformVelocity;
};

#endif //MAIN_H
