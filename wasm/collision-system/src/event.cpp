//
// Created by Sheldon Lobo on 10/14/24.
//

#include <iostream>
#include <iomanip>
#include <cassert>

#include "event.h"
#include "ball.h"

Event::Event() :
    time(0.0), a(nullptr), b(nullptr), countA(-1), countB(-1) {}

Event::Event(double time, Ball *a, Ball *b) :
    time(time), a(a), b(b) {
    // std::cout << "Event: " << *a << " <-> " << *b << " t: " << std::fixed << std::setprecision(2) << time << std::endl;
    assert(time >= 0.0);

    if (a != nullptr) {
        countA = a->getCount();
    } else {
        countA = -1;
    }
    if (b != nullptr) {
        countB = b->getCount();
    } else {
        countB = -1;
    }

    // std::cout << "Event::Event : time: " << std::fixed << std::setprecision(2) << time << " ";
    // if (a != nullptr) {
    //     std::cout << a << " ";
    // } else {
    //     std::cout << "[null] ";
    // }
    // if (b != nullptr) {
    //     std::cout << b << "\n";
    // } else {
    //     std::cout << "[null]\n";
    // }
}

bool Event::operator<(const Event& e) const {
    return time > e.time;
}

bool Event::isValid() const {
    if (a != nullptr && a->getCount() != countA) {
        return false;
    }
    if (b != nullptr && b->getCount() != countB) {
        return false;
    }

    return true;
}
