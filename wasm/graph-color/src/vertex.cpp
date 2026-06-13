//
// Created by Sheldon Lobo on 5/29/26.
//

#include "vertex.h"

#include <cassert>

#include "main.h"
#include "sdl-helpers/circle.h"

#include <iostream>

Vertex::Vertex(unsigned x, unsigned y, unsigned id) : id(id), x(x), y(y), color(-1) {
    euclideanDistance = x*x + y*y;
}

double Vertex::distanceTo(Vertex* v) const {
    int xDiff = x - v->x;
    int yDiff = y - v->y;
    return std::sqrt(xDiff * xDiff + yDiff * yDiff);
}

bool Vertex::inRange(int givenX, int givenY) const {
    return (givenX >= (x - IN_RANGE) &&  givenX <= (x + IN_RANGE) && givenY >= (y - IN_RANGE) &&
        givenY <= y + IN_RANGE);
}

bool Vertex::tooClose(int givenX, int givenY) const {
    return (givenX >= (x - SEPARATION) &&  givenX <= (x + SEPARATION) && givenY >= (y - SEPARATION) &&
        givenY <= y + SEPARATION);
}

void Vertex::removeEdge(const Edge* edge) {
    for (unsigned i = 0; i < edges.size(); i++) {
        if (edges[i] == edge) {
            // std::cout << "    Removing edge: " << *edges[i] << std::endl;
            edges.erase(edges.begin() + i);
            break;
        }
    }
}

int Vertex::degree() const {
    return edges.size();
}

SDL_Color Vertex::drawColor() const {
    SDL_Color sdlColor;
    switch (color) {
        case 0:
            sdlColor = SDL_Color{100, 255, 100, SDL_ALPHA_OPAQUE};
            break;
        case 1:
            sdlColor = SDL_Color{255, 100, 100, SDL_ALPHA_OPAQUE};
            break;
        case 2:
            sdlColor = SDL_Color{100, 100, 255, SDL_ALPHA_OPAQUE};
            break;
        case 3:
            sdlColor = SDL_Color{255, 255, 50, SDL_ALPHA_OPAQUE};
            break;
        default:
            assert(false);
    }
    return sdlColor;
}

void Vertex::draw(Context* ctx) const {
    DrawFilledCircle(ctx->renderer, x, y, ctx->vertexRadius, drawColor());

    // Write id in circle
    SDL_Color black = {0, 0, 0};
    std::string idString = std::to_string(id);
    SDL_Surface* surfaceMessage = TTF_RenderText_Solid(ctx->font, idString.c_str(), black);
    SDL_Texture* Message = SDL_CreateTextureFromSurface(ctx->renderer, surfaceMessage);

    SDL_Rect Message_rect; //create a rect
    if (id <= 9) {
        Message_rect.x = x - RADIUS/2;
        Message_rect.w = RADIUS;
    } else {
        Message_rect.x = x - RADIUS;
        Message_rect.w = RADIUS*2;
    }
    Message_rect.y = y - RADIUS;
    Message_rect.h = RADIUS*2;

    SDL_RenderCopy(ctx->renderer, Message, NULL, &Message_rect);

    SDL_FreeSurface(surfaceMessage);
    SDL_DestroyTexture(Message);
}

bool Vertex::operator<(const Vertex& other) const {
    return id < other.id;
}

std::ostream& operator<<(std::ostream& strm, const Vertex& v) {
    strm << "{" << v.id << "} [(" << v.x << "," << v.y << "); " << v.color << " Edges:" << v.edges.size() << "]";
    return strm;
}

bool EuclideanDistanceComparator::operator()(const std::unique_ptr<Vertex>& v1, const std::unique_ptr<Vertex>& v2) const {
    return v1->euclideanDistance < v2->euclideanDistance;
}
