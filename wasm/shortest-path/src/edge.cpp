//
// Created by Sheldon Lobo on 4/25/25.
//

#include "edge.h"
#include "main.h"
#include "vertex.h"

#include <iomanip>

DirectedEdge::DirectedEdge(Vertex *from, Vertex *to) :
    from(from), to(to) {
    // weight = std::sqrt(std::pow(from->x - to->x, 2) + std::pow(from->y - to->y, 2));

    // Calculate the arrow x & y coordinates
    int xDiff = from->x - to->x;
    int yDiff = from->y - to->y;

    weight = std::sqrt(xDiff * xDiff + yDiff * yDiff);

    int xOffset = RADIUS/weight * xDiff;
    int yOffset = RADIUS/weight * yDiff;

    xFrom = from->x - xOffset;
    yFrom = from->y - yOffset;
    xTo = to->x + xOffset;
    yTo = to->y + yOffset;

    // Calculate the angle
    angle = std::atan2(yTo - yFrom, xTo - xFrom);
    angleDegrees = angle * 180.0 / M_PI;
    double leftArrowAngle = angle - M_PI/6;
    double rightArrowAngle = angle + M_PI/6;
    xLeft = xTo - ARROW * std::cos(leftArrowAngle);
    yLeft = yTo - ARROW * std::sin(leftArrowAngle);
    xRight = xTo - ARROW * std::cos(rightArrowAngle);
    yRight = yTo - ARROW * std::sin(rightArrowAngle);
}

void DirectedEdge::draw(Context *ctx, SDL_Color color) {
    SDL_SetRenderDrawColor(ctx->renderer, color.r, color.g, color.b, color.a);
    SDL_RenderDrawLine(ctx->renderer, xFrom, yFrom, xTo, yTo);
    SDL_RenderDrawLine(ctx->renderer, xLeft, yLeft, xTo, yTo);
    SDL_RenderDrawLine(ctx->renderer, xRight, yRight, xTo, yTo);
}

std::ostream& operator<<(std::ostream &strm, const DirectedEdge &e) {
    strm << std::fixed << std::setprecision(4);
    strm << "<" << *e.from << " -> " << *e.to << " (" << e.weight << "; " << e.angleDegrees << ")" << ">";
    return strm;
}

bool EdgeWeightComparator::operator()(const DirectedEdge *e1, const DirectedEdge *e2) const {
    return e1->weight < e2->weight;
}

bool EdgeFromComparator::operator()(const DirectedEdge *e1, const DirectedEdge *e2) const {
    if (e1->from->id != e2->from->id) {
        return e1->from->id < e2->from->id;
    }
    return e1->to->id < e2->to->id;
}
