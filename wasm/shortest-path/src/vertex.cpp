//
// Created by Sheldon Lobo on 4/25/25.
//

#include "vertex.h"
#include "main.h"
#include "sdl-helpers/circle.h"

#include <iostream>

Vertex::Vertex(unsigned x, unsigned y, unsigned id, unsigned origId) : id(id), origId(origId), x(x), y(y) {
    euclideanDistance = x*x + y*y;
}

void Vertex::setId(unsigned id) {
    this->id = id;
    this->origId = id;
}

double Vertex::distanceTo(Vertex *v) const {
    int xDiff = x - v->x;
    int yDiff = y - v->y;
    return std::sqrt(xDiff * xDiff + yDiff * yDiff);
}

bool Vertex::inRange(int givenX, int givenY) const {
    return (givenX >= (x - SEPARATION) &&  givenX <= (x + SEPARATION) && givenY >= (y - SEPARATION) &&
        givenY <= y + SEPARATION);
}

void Vertex::removeIncomingEdge(const DirectedEdge *edge) {
    for (unsigned i = 0; i < edgesTo.size(); i++) {
        if (edgesTo[i].get() == edge) {
            std::cout << "    Removing from successor: " << *edgesTo[i] << std::endl;
            edgesTo.erase(edgesTo.begin() + i);
            break;
        }
    }
}

void Vertex::removeOutgoingEdge(const DirectedEdge *edge) {
    for (unsigned i = 0; i < edgesFrom.size(); i++) {
        if (edgesFrom[i].get() == edge) {
            std::cout << "    Removing from predecessor: " << *edgesFrom[i] << std::endl;
            edgesFrom.erase(edgesFrom.begin() + i);
            break;
        }
    }
}

void Vertex::draw(Context *ctx, SDL_Color color) {
    DrawFilledCircle(ctx->renderer, x, y, ctx->vertexRadius, color);

    // Write id in circle
    SDL_Color black = {0, 0, 0};
    std::string idString = std::to_string(origId);
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

std::ostream& operator<<(std::ostream &strm, const Vertex &v) {
    strm << "{" << v.id << "/" << v.origId << "} [" << v.x << "," << v.y << "; Froms:" << v.edgesFrom.size() <<
        ", Tos:" << v.edgesTo.size() << "]";
    return strm;
}

bool EuclideanDistanceComparator::operator()(const std::unique_ptr<Vertex>& v1, const std::unique_ptr<Vertex>& v2) const {
    return v1->euclideanDistance < v2->euclideanDistance;
}
