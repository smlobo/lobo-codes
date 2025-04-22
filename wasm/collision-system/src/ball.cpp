//
// Created by Sheldon Lobo on 10/6/24.
//

#include <limits>
#include <cassert>
#include <cmath>
#include <iomanip>

#include "ball.h"

#include <iostream>

#include "sdl-helpers/circle.h"

Ball::Ball(double x, double y, double vx, double vy, double radius, double mass) :
    x(x), y(y), vx(vx), vy(vy), radius(radius), mass(mass), count(0) {}

void Ball::move(double dt) {
    x += vx * dt;
    y += vy * dt;
    // std::cout << "Ball::move : " << *this << std::endl;
}

double Ball::timeToHit(const Ball *that) const {
    if (this == that) {
        return std::numeric_limits<double>::infinity();
    }

    double dx = that->x - this->x;
    double dy = that->y - this->y;
    double sigma = this->radius + that->radius;

    // Already overlapping
    if (std::abs(dx) < sigma && std::abs(dy) < sigma) {
        return std::numeric_limits<double>::infinity();
    }

    double dvx = that->vx - this->vx;
    double dvy = that->vy - this->vy;
    double dvdr = dx*dvx + dy*dvy;

    if (dvdr > 0) {
        return std::numeric_limits<double>::infinity();
    }

    double dvdv = dvx*dvx + dvy*dvy;
    double drdr = dx*dx + dy*dy;
    double d = (dvdr*dvdr) - dvdv * (drdr - sigma*sigma);
    if (d < 0) {
        return std::numeric_limits<double>::infinity();
    }

    double rv = -(dvdr + std::sqrt(d)) / dvdv;
    // std::cout << *this << " <-> " << *that << " = " << std::fixed << std::setprecision(2) << rv << std::endl;
    assert (rv >= 0.0);
    return rv;
}

double Ball::timeToHitVertical() const {
    // Always return a +ve time
    double time = std::numeric_limits<double>::infinity();
    if (vx > 0.0) {
        time = (1.0 - x - radius) / vx;
    } else if (vx < 0.0) {
        time = (0.0 - (x - radius)) / vx;
    }
    return std::abs(time);
}

double Ball::timeToHitHorizontal() const {
    // Always return a +ve time
    double time = std::numeric_limits<double>::infinity();
    if (vy > 0.0) {
        time = (1.0 - y - radius) / vy;
    } else if (vy < 0.0) {
        time = (0.0 - (y - radius)) / vy;
    }
    return std::abs(time);
}

void Ball::bounceOff(Ball *that) {
    double dx  = that->x - this->x;
    double dy  = that->y - this->y;
    double dvx = that->vx - this->vx;
    double dvy = that->vy - this->vy;
    double dvdr = dx*dvx + dy*dvy;
    double dist = this->radius + that->radius;

    double J = 2 * this->mass * that->mass * dvdr /
        ((this->mass + that->mass) * dist);
    double Jx = J * dx / dist;
    double Jy = J * dy / dist;

    this->vx += Jx / this->mass;
    this->vy += Jy / this->mass;
    that->vx -= Jx / that->mass;
    that->vy -= Jy / that->mass;
    this->count++;
    that->count++;
}

void Ball::bounceOffVertical() {
    vx = -vx;
    count++;
}

void Ball::bounceOffHorizontal() {
    vy = -vy;
    count++;
}

unsigned Ball::getCount() const {
    return count;
}

SDL_Color Ball::getColor() const {
    unsigned char red = 0, green = red, blue = red;

    unsigned char step = count % 25 * 10;

    unsigned trimester = (count / 25) % 3;
    switch (trimester) {
        case 0:
            blue = 250 - step;
            green = step;
            break;
        case 1:
            green = 250 - step;
            red = step;
            break;
        case 2:
            red = 250 - step;
            blue = step;
            break;
    }

    return SDL_Color{red, green, blue, SDL_ALPHA_OPAQUE};
}

void Ball::draw(Context *context) const {
    int x = this->x * context->xDimension;
    int y = this->y * context->yDimension;
    int r = this->radius * context->scale;
    DrawFilledCircle(context->renderer, x, y, r, getColor());
}

std::ostream& operator<<(std::ostream& strm, const Ball& b) {
    strm << std::fixed << std::setprecision(4);
    strm << "[" << b.x << "," << b.y << "," <<  b.vx << "," << b.vy << "," << b.count << "]";
    return strm;
}
