#include <stringapiset.h>
#include "ImageScale.h"
#include <opencv2/opencv.hpp>

using namespace cv;

int ImageScale::scale(char *srcPath, char *savePath, int width, int height)
{
    const char *gbkSrcPath = utf8toGbk(srcPath);
    const char *gbkSavePath = utf8toGbk(savePath);
    Mat src = imread(gbkSrcPath);
    if (!src.data)
    {
        return -1;
    }
    Size srcSize = src.size();
    int interpolation = INTER_AREA;
    if (width > srcSize.width || height > srcSize.height)
    {
        interpolation = INTER_CUBIC;
    }
    Size saveSize(width, height);
    Mat save;
    resize(src, save, saveSize, interpolation);
    int result = imwrite(gbkSavePath, save) ? 0 : -1;
    src.release();
    save.release();
    return result;
}

int ImageScale::scale(char *srcPath, char *savePath, double percentage)
{
    const char *gbkSrcPath = utf8toGbk(srcPath);
    const char *gbkSavePath = utf8toGbk(savePath);

    Mat src = imread(gbkSrcPath);
    if (!src.data)
    {
        return -1;
    }
    Size srcSize = src.size();
    int width = (int)(srcSize.width * percentage);
    int height = (int)(srcSize.height * percentage);
    int interpolation = INTER_AREA;
    if (percentage > 1)
    {
        interpolation = INTER_CUBIC;
    }
    Size saveSize(width, height);
    Mat save;
    resize(src, save, saveSize, interpolation);
    int result = imwrite(gbkSavePath, save) ? 0 : -1;
    src.release();
    save.release();
    return result;
}

char *ImageScale::utf8toGbk(const char *str)
{
    int size;
    size = MultiByteToWideChar(CP_UTF8, 0, str, -1, nullptr, 0);
    auto *strUnicode = new wchar_t[size];
    MultiByteToWideChar(CP_UTF8, 0, str, -1, strUnicode, size);
    size = WideCharToMultiByte(CP_ACP, 0, strUnicode, -1, nullptr, 0, nullptr, nullptr);
    char *strGbk = new char[size];
    WideCharToMultiByte(CP_ACP, 0, strUnicode, -1, strGbk, size, nullptr, nullptr);
    delete[] strUnicode;
    return strGbk;
}
