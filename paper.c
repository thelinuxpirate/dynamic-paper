#!/usr/bin/env scriptisto

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

// scriptisto-begin
// script_src: paper.c
// build_cmd: clang -O2 paper.c -o ./script
// scriptisto-end


// current time using tm from time.h
struct tm currentTime() {
    time_t rawtime;
    struct tm * timeinfo;
    time(&rawtime);
    timeinfo = localtime(&rawtime); 
    return *timeinfo;
}

// executes OS command
void execCmd(const char* cmd) {
    system(cmd);
}

// using feh set wallpaper depending on the time
void determineTime() {
    const int sunriseHour = 6;
    const int dayHour = 11;
    const int sunsetHour = 19;
    const int nightHour = 20;

    struct tm now = currentTime();

    if (now.tm_hour >= nightHour) {
        execCmd("feh --bg-fill ~/Pictures/Wallpapers/etc/outset-island/night.jpg");
        printf("Wallpaper set to: Nighttime\n");
    } else if (now.tm_hour >= sunsetHour) {
        execCmd("feh --bg-fill ~/Pictures/Wallpapers/etc/outset-island/beforeYafter.png");
        printf("Wallpaper set to: Sunset\n");
    } else if (now.tm_hour >= dayHour) {
        execCmd("feh --bg-fill ~/Pictures/Wallpapers/etc/outset-island/day.png");
        printf("Wallpaper set to: Daytime\n");
    } else if (now.tm_hour >= sunriseHour) {
        execCmd("feh --bg-fill ~/Pictures/Wallpapers/etc/outset-island/beforeYafter.png");
        printf("Wallpaper set to: Sunrise\n");
    }
}

int main() {
    determineTime();

    struct tm now = currentTime();
    printf("Current hour: %d\n", now.tm_hour); // print hour

    return 0;
}

