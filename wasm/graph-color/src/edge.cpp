//
// Created by Sheldon Lobo on 5/29/26.
//

#include "edge.h"
#include "main.h"
#include "vertex.h"

#include <iomanip>

Edge::Edge(Vertex *from, Vertex *to) :
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
}

bool Edge::intersects(Edge *edge) {
    if (edge == nullptr) {
        return false;
    }

    // Edges that share a vertex are allowed to meet at that endpoint.
    if (from == edge->from || from == edge->to || to == edge->from || to == edge->to) {
        return false;
    }

    auto orientation = [](int ax, int ay, int bx, int by, int cx, int cy) {
        long long cross = static_cast<long long>(bx - ax) * (cy - ay) -
            static_cast<long long>(by - ay) * (cx - ax);
        if (cross == 0) {
            return 0;
        }
        return cross > 0 ? 1 : -1;
    };

    auto onSegment = [](int ax, int ay, int bx, int by, int px, int py) {
        return px >= std::min(ax, bx) && px <= std::max(ax, bx) &&
            py >= std::min(ay, by) && py <= std::max(ay, by);
    };

    int o1 = orientation(xFrom, yFrom, xTo, yTo, edge->xFrom, edge->yFrom);
    int o2 = orientation(xFrom, yFrom, xTo, yTo, edge->xTo, edge->yTo);
    int o3 = orientation(edge->xFrom, edge->yFrom, edge->xTo, edge->yTo, xFrom, yFrom);
    int o4 = orientation(edge->xFrom, edge->yFrom, edge->xTo, edge->yTo, xTo, yTo);

    if (o1 != o2 && o3 != o4) {
        return true;
    }

    if (o1 == 0 && onSegment(xFrom, yFrom, xTo, yTo, edge->xFrom, edge->yFrom)) {
        return true;
    }
    if (o2 == 0 && onSegment(xFrom, yFrom, xTo, yTo, edge->xTo, edge->yTo)) {
        return true;
    }
    if (o3 == 0 && onSegment(edge->xFrom, edge->yFrom, edge->xTo, edge->yTo, xFrom, yFrom)) {
        return true;
    }
    if (o4 == 0 && onSegment(edge->xFrom, edge->yFrom, edge->xTo, edge->yTo, xTo, yTo)) {
        return true;
    }

    return false;
}

bool Edge::intersects2(const Edge *edge) {
    static constexpr double eps = 1e-9;

    // Edges that share a vertex are allowed to meet at that endpoint.
    if (from == edge->from || from == edge->to || to == edge->from || to == edge->to) {
        return false;
    }

    double det = (from->x - to->x) * (edge->from->y - edge->to->y) - (from->y - to->y) * (edge->from->x - edge->to->x);

    if (std::abs(det) < eps) {
        return false; // parallel or collinear
    }

    double t1 = from->x * to->y - from->y * to->x;
    double t2 = edge->from->x * edge->to->y - edge->from->y * edge->to->x;

    double intersectX = (t1 * (edge->from->x - edge->to->x) - (from->x - to->x) * t2) / det;
    double intersectY = (t1 * (edge->from->y - edge->to->y) - (from->y - to->y) * t2) / det;

    auto onSegment = [](const Edge *e, double ix, double iy) -> bool {
        return ix >= std::min(e->from->x, e->to->x) - eps && ix <= std::max(e->from->x, e->to->x) + eps &&
               iy >= std::min(e->from->y, e->to->y) - eps && iy <= std::max(e->from->y, e->to->y) + eps;
    };
    return onSegment(this, intersectX, intersectY) && onSegment(edge, intersectX, intersectY);
}

void Edge::draw(Context *ctx, SDL_Color color) const {
    SDL_SetRenderDrawColor(ctx->renderer, color.r, color.g, color.b, color.a);
    SDL_RenderDrawLine(ctx->renderer, xFrom, yFrom, xTo, yTo);
}

bool Edge::operator==(const Edge& other) const {
    return from->id == other.from->id &&
           to->id == other.to->id;
}

std::ostream& operator<<(std::ostream &strm, const Edge &e) {
    strm << std::fixed << std::setprecision(4);
    strm << "<" << *e.from << " -> " << *e.to << " (" << e.weight << "; " << e.angleDegrees << ")" << ">";
    return strm;
}

std::size_t EdgeHash::operator()(const Edge& e) const {
    return std::hash<unsigned>{}(e.from->id) ^ (std::hash<unsigned>{}(e.to->id) << 1);
}

bool EdgeWeightComparator::operator()(const Edge &e1, const Edge &e2) const {
    return e1.weight < e2.weight;
}

bool EdgeFromComparator::operator()(const Edge &e1, const Edge &e2) const {
    if (e1.from->id != e2.from->id) {
        return e1.from->id < e2.from->id;
    }
    return e1.to->id < e2.to->id;
}
