# image-resize
使用go实现的命令行批量图片大小调整工具

### 使用
---

下载 [image-resize_win64.zip](https://github.com/Tait4198/image-resize/releases/download/v1.0/image-resize_win64.zip) 后解压可直接使用

**使用示例**

示例1
>image-resize.exe -i C:\image\a.jpg -s 1920x1080
>
>将 `C:\image\a.jpg` 调整为宽 1920px 高 1080px	输出文件` C:\image\a_1920x1080.jpg`
>

示例2
>image-resize.exe -i C:\image -p jpg -e png -s 0.25 -d 1 
>将 `C:\image` 文件夹第1层下(不包含子文件夹)的jpg格式图片缩放为 `25%`并且存储为png格式图片

示例3

> image-resize.exe -i C:\image -o C:\new_image  -s 2 -g 4
>
> 将 `C:\image` 文件夹下所有默认类型图片放大为 `200%`并且存储至 `C:\new_image` 文件夹下
>
> `C:\image\a.jpg` 将会被存储至 `C:\new_image\a_200%.jpg`
>
> `C:\image\jpg\b.jpg` 将会被存储至 `C:\new_image\jpg\b_200%.jpg`



### 参数
---

| 参数 | 说明                                                         | 默认         |
| :--: | ------------------------------------------------------------ | ------------ |
|  i   | 原始图像或图像文件夹位置<br />例. C:\\a.jpg C:\\image        |              |
|  o   | 处理后图像存储文件夹位置（默认存储于原始图所在文件夹）       |              |
|  p   | 将被处理的图像类型（多个类型使用 / 隔开)                     | jpg/jpeg/png |
|  e   | 处理后输出的图像类型（需要使用有效的图像类型）<br />默认使用原始文件类型 |              |
|  s   | 调整图像尺寸<br />例. 百分比：0.5 调整为50% 指定宽高：200x300 调整为宽200高300 |              |
|  d   | 对子文件夹处理的层数（0表示对所有子文件夹处理,1表示当前文件夹,2表示包含一层子文件夹以此类推） | 0            |
|  g   | 并行处理数量（默认使用逻辑CPU数量）                          |              |