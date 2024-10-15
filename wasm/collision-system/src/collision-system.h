//
// Created by Sheldon Lobo on 10/12/24.
//

#ifndef COLLISION_SYSTEM_H
#define COLLISION_SYSTEM_H

#include "main.h"

void CollisionSystemInit(Context &);
void CollisionSystemSimulate(Context *);
void CollisionSystemAdd(Context *);
void CollisionSystemSubtract(Context *);

#endif //COLLISION_SYSTEM_H
