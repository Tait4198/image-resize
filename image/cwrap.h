#ifndef IMAGE_SCALE_CWRAP_H
#define IMAGE_SCALE_CWRAP_H

#ifdef __cplusplus
extern "C"
{
#endif

    int scaleSize(char *, char *, int, int);

    int scalePercentage(char *, char *, double);

#ifdef __cplusplus
};
#endif

#endif //IMAGE_SCALE_CWRAP_H
