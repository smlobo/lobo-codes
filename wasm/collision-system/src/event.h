//
// Created by Sheldon Lobo on 10/14/24.
//

#ifndef EVENT_H
#define EVENT_H

class Ball;

class Event {
private:
    double time;
    Ball *a;
    Ball *b;
    int countA;
    int countB;

public:
    Event();
    Event(double time, Ball *a, Ball *b);

    bool operator<(const Event& e) const;
    bool isValid() const;

    double getTime() const { return time; };
    Ball *getA() const { return a; };
    Ball *getB() const { return b; };
};

#endif //EVENT_H
