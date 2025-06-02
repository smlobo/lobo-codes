//
// Created by Sheldon Lobo on 4/25/25.
//

#include "vertex.h"
#include "main.h"
#include "sdl-helpers/circle.h"

#include <iostream>

Vertex::Vertex(unsigned x, unsigned y) : id(0), x(x), y(y) {
    euclideanDistance = x*x + y*y;
}

double Vertex::distanceTo(Vertex &v) {
    int xDiff = x - v.x;
    int yDiff = y - v.y;
    return std::sqrt(xDiff * xDiff + yDiff * yDiff);
}

void Vertex::draw(Context *ctx, SDL_Color color) {
    DrawFilledCircle(ctx->renderer, x, y, ctx->vertexRadius, color);

    // Write id in circle
    SDL_Color black = {0, 0, 0};
    std::string idString = std::to_string(id);
    SDL_Surface* surfaceMessage = TTF_RenderText_Solid(ctx->font, idString.c_str(), black);
    SDL_Texture* Message = SDL_CreateTextureFromSurface(ctx->renderer, surfaceMessage);

    SDL_Rect Message_rect; //create a rect
    Message_rect.x = x - RADIUS/2;  //controls the rect's x coordinate
    Message_rect.y = y - RADIUS/2; // controls the rect's y coordinte
    Message_rect.w = RADIUS; // controls the width of the rect
    Message_rect.h = RADIUS; // controls the height of the rect

    SDL_RenderCopy(ctx->renderer, Message, NULL, &Message_rect);

    SDL_FreeSurface(surfaceMessage);
    SDL_DestroyTexture(Message);
}

std::ostream& operator<<(std::ostream &strm, const Vertex &v) {
    strm << "{" << v.id << "} [" << v.x << "," << v.y << "; Froms:" << v.edgesFrom.size() << ", Tos:" <<
        v.edgesTo.size() << "]";
    return strm;
}

bool EuclideanDistanceComparator::operator()(const Vertex &v1, const Vertex &v2) const {
    return v1.euclideanDistance < v2.euclideanDistance;
}
