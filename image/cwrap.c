#include "cwrap.h"
#include "ImageScale.h"

#define Ir ImageScale
Ir ir;

int scaleSize(char *src, char *save, int width, int height)
{
    return ir.scale(src, save, width, height);
}

int scalePercentage(char *src, char *save, double percentage)
{
    return ir.scale(src, save, percentage);
}