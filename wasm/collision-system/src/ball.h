//
// Created by Sheldon Lobo on 10/6/24.
//

#ifndef BALL_H
#define BALL_H
#include <SDL_pixels.h>

#include "main.h"

class Ball {
private:
    double x, y;
    double vx, vy;
    double const radius;
    double const mass;
    unsigned count;

public:
    Ball(double x, double y, double vx, double vy, double radius, double mass);

    void move(double dt);

    double timeToHit(const Ball *that) const;
    double timeToHitVertical() const;
    double timeToHitHorizontal() const;

    void bounceOff(Ball *that);
    void bounceOffVertical();
    void bounceOffHorizontal();

    unsigned getCount() const;
    SDL_Color getColor() const;

    void draw(Context *context) const;

    friend std::ostream& operator<<(std::ostream&, const Ball&);
};

#endif //BALL_H
