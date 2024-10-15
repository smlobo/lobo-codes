//
// Created by Sheldon Lobo on 10/6/24.
//

#include <cassert>
#include <iostream>
#include <random>

#include "ball.h"
#include "event.h"
#include "main.h"

constexpr double RADIUS = 0.01;
constexpr double MASS = 0.5;

void redraw(Context *ctx) {
    // Clear
    SDL_SetRenderDrawColor(ctx->renderer, 255, 255, 255, SDL_ALPHA_OPAQUE);
    SDL_RenderClear(ctx->renderer);

    // Draw
    for (const Ball *ball : ctx->balls) {
        ball->draw(ctx);
    }

    // Present
    SDL_RenderPresent(ctx->renderer);

    ctx->events.emplace(ctx->time + 2.0, nullptr, nullptr);
}

void predict(Context *ctx, Ball *a) {
    if (a == nullptr)
        return;

    // std::cout << "predict : " << *a << std::endl;
    for (Ball *ball : ctx->balls) {
        double dt = a->timeToHit(ball);
        ctx->events.emplace(ctx->time + dt, a, ball);
    }
    ctx->events.emplace(ctx->time + a->timeToHitVertical(), a, nullptr);
    ctx->events.emplace(ctx->time + a->timeToHitHorizontal(), nullptr, a);
}

void CollisionSystemInit(Context &ctx) {
    // Create balls
    for (int i = 0; i < ctx.nBalls; i++) {
        Ball *ball = new Ball(ctx.uniformLocation(ctx.re), ctx.uniformLocation(ctx.re),
            ctx.uniformVelocity(ctx.re), ctx.uniformVelocity(ctx.re), RADIUS, MASS);
        ctx.balls.push_back(ball);
    }

    // Initial predictions
    for (Ball *ball : ctx.balls) {
        predict(&ctx, ball);
    }
    // Create initial event
    ctx.events.emplace(ctx.time, nullptr, nullptr);
}

void CollisionSystemAdd(Context *ctx) {
    // Add a ball
    Ball *ball = new Ball(ctx->uniformLocation(ctx->re), ctx->uniformLocation(ctx->re),
        ctx->uniformVelocity(ctx->re), ctx->uniformVelocity(ctx->re), RADIUS, MASS);
    ctx->balls.push_back(ball);

    // Clear events
    ctx->events = std::priority_queue<Event>();

    // Initial predictions
    for (Ball *b : ctx->balls) {
        predict(ctx, b);
    }

    redraw(ctx);
}

void CollisionSystemSubtract(Context *ctx) {
    // Remove the first ball
    ctx->balls.pop_front();

    // Clear events
    ctx->events = std::priority_queue<Event>();

    // Initial predictions
    for (Ball *b : ctx->balls) {
        predict(ctx, b);
    }

    redraw(ctx);
}

void CollisionSystemSimulate(Context *ctx) {
    // Handle all Events until the redraw event
    while (!ctx->events.empty()) {
        Event event = ctx->events.top();
        ctx->events.pop();
        if (!event.isValid()) {
            continue;
        }

        // std::cout << "CollisionSystemSimulate : [" << ctx->events.size() << "]" << std::endl;
        Ball *a = event.getA();
        Ball *b = event.getB();
        for(Ball *ball: ctx->balls) {
            ball->move(event.getTime() - ctx->time);
        }
        ctx->time = event.getTime();

        if (a != nullptr && b != nullptr) {
            a->bounceOff(b);
        } else if (a != nullptr) {
            a->bounceOffVertical();
        } else if (b != nullptr) {
            b->bounceOffHorizontal();
        }

        predict(ctx, a);
        predict(ctx, b);

        if (a == nullptr && b == nullptr) {
            redraw(ctx);
            break;
        }
    }
}