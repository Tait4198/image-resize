#ifndef IMAGE_SCALE_IMAGESCALE_H
#define IMAGE_SCALE_IMAGESCALE_H

class ImageScale
{
public:
    static int scale(char *, char *, int, int);
    static int scale(char *, char *, double);

private:
    static char *utf8toGbk(const char *);
};

#endif //IMAGE_SCALE_IMAGESCALE_H
