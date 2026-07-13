//
// Created by Sheldon Lobo on 7/7/26.
//

#include "cell.h"

#include <functional>

#include "main.h"

constexpr Cell NULL_CELL{-1, -1};

int Cell::midX() const {
    return x*DIAMETER + RADIUS;
}

int Cell::midY() const {
    return y*DIAMETER + RADIUS;
}

Cell cellFor(int x, int y) {
    return {x/DIAMETER, y/DIAMETER};
}

std::size_t CellHash::operator()(const Cell& cell) const {
    return std::hash<int>()(cell.x) ^ (std::hash<int>()(cell.y) << 1);
}
