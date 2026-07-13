//
// Created by Sheldon Lobo on 7/7/26.
//

#ifndef MINIMUM_SPANNING_TREE_CELL_H
#define MINIMUM_SPANNING_TREE_CELL_H

#include <cstddef>

class Vertex;

struct Cell {
    int x;
    int y;

    int midX() const;
    int midY() const;

    bool operator==(const Cell& other) const = default;
};

extern const Cell NULL_CELL;

Cell cellFor(int x, int y);

struct CellHash {
    std::size_t operator()(const Cell& cell) const;
};

#endif //MINIMUM_SPANNING_TREE_CELL_H