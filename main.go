package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

func main() {
	var srcPath, savePath, inFormat, outFormat, size string
	var depth, gs int
	required := []string{"i", "s"}

	flag.StringVar(&srcPath, "i", "", "原始图像或图像文件夹位置\n例. C:\\a.jpg C:\\image")
	flag.StringVar(&savePath, "o", "", "处理后图像存储文件夹位置（默认存储于原始图所在文件夹）")
	flag.StringVar(&inFormat, "p", "jpg/jpeg/png", "将被处理的图像类型（多个类型使用 / 隔开")
	flag.StringVar(&outFormat, "e", "", "处理后输出的图像类型（需要使用有效的图像类型）\n默认根据使用原始文件类型")
	flag.StringVar(&size, "s", "", "调整图像尺寸\n例. 百分比：0.5 调整为50% 指定宽高：200x300 调整为宽200高300")
	flag.IntVar(&depth, "d", 0, "对子文件夹处理的层数（0表示对所有子文件夹处理,1表示当前文件夹,2表示包含一层子文件夹以此类推）")
	flag.IntVar(&gs, "g", 0, "并行处理数量（默认使用逻辑CPU数量）")
	flag.Parse()

	seen := make(map[string]bool)
	valid := true
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			valid = false
			fmt.Printf("缺少必要参数 -%s\n", req)
		}
	}
	if !valid {
		fmt.Printf("使用 -h 查看命令介绍\n")
		os.Exit(-1)
	}
	var rs ResizeSize
	if match, _ := regexp.MatchString("^[\\d+]*([.]\\d+)?$", size); match {
		p, _ := strconv.ParseFloat(size, 64)
		rs = ResizeSize{Percentage: p, Width: 0, Height: 0}
	} else if match, _ := regexp.MatchString("^\\d+[xX]\\d+$", size); match {
		compile, _ := regexp.Compile("(\\d+)[xX](\\d+)")
		m := compile.FindSubmatch([]byte(size))
		width, _ := strconv.Atoi(string(m[1]))
		height, _ := strconv.Atoi(string(m[2]))
		rs = ResizeSize{Percentage: 0, Width: width, Height: height}
	} else {
		log.Printf("调整图像尺寸参数无效\n")
		os.Exit(-1)
	}
	if gs == 0 {
		gs = runtime.NumCPU()
	}

	t1 := time.Now()
	err := Resize(srcPath, savePath, inFormat, outFormat, rs, depth, gs)
	if err != nil {
		os.Exit(-1)
	}
	t2 := time.Since(t1)
	log.Printf("处理完成\t总耗时:%s", t2)
	os.Exit(0)
}
