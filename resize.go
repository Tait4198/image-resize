package main

// #cgo CFLAGS: -Iimage
// #cgo LDFLAGS: -Limage -limage -lstdc++
// #cgo LDFLAGS: -L. -lopencv_core411 -lopencv_imgcodecs411 -lopencv_imgproc411 -lstdc++
// #include "image/cwrap.h"
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

type ResizeFile struct {
	SrcPath  string
	SavePath string
}

func (rs ResizeSize) String() string {
	if rs.Percentage > 0 {
		return fmt.Sprintf("%.0f%%", rs.Percentage*100)
	} else {
		return fmt.Sprintf("%dx%d", rs.Width, rs.Height)
	}
}

type ResizeSize struct {
	Percentage float64
	Width      int
	Height     int
}

type ResizeFileFilter func(string) bool

type FileNameConversion func(string) string

func Resize(srcPath, savePath, inFormat, outFormat string, size ResizeSize, depth, gs int) error {
	formatMap := make(map[string]bool)
	for _, format := range strings.Split(inFormat, "/") {
		formatMap["."+format] = true
	}
	rff := func(fileName string) bool {
		startIndex := strings.Index(fileName, ".")
		_, ex := formatMap[fileName[startIndex:]]
		return !ex
	}
	fnc := func(fileName string) string {
		startIndex := strings.Index(fileName, ".")
		format := fileName[startIndex:]
		name := fileName[:startIndex]
		if outFormat != "" {
			return fmt.Sprintf("%s_%s.%s", name, size, outFormat)
		} else {
			return fmt.Sprintf("%s_%s%s", name, size, format)
		}
	}
	resizeFiles, err := GetResizeFiles(srcPath, savePath, "", 0, depth, rff, fnc)
	if err != nil {
		return err
	}
	if savePath != "" {
		if err := MkDirs(savePath); err != nil {
			return err
		}
	}

	taskChan := make(chan ResizeFile, gs)
	exitChan := make(chan bool, gs)
	go func() {
		for _, rf := range resizeFiles {
			taskChan <- rf
		}
		close(taskChan)
	}()
	for i := 0; i < gs; i++ {
		go ResizeTask(taskChan, exitChan, size)
	}
	for i := 0; i < gs; i++ {
		<-exitChan
	}
	close(exitChan)
	return nil
}

func ResizeTask(taskChan chan ResizeFile, exitChan chan bool, size ResizeSize) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("task err %s\n", err)
			return
		}
	}()

	for task := range taskChan {
		fmt.Println(task.SrcPath, task.SavePath)
		ResizeImage(task, size)
	}
	exitChan <- true
}

func ResizeImage(rf ResizeFile, rs ResizeSize) {
	src := C.CString(rf.SrcPath)
	save := C.CString(rf.SavePath)
	if rs.Percentage == 0 {
		_ = C.scaleSize(src, save, C.int(rs.Width), C.int(rs.Height))
	} else {
		_ = C.scalePercentage(src, save, C.double(rs.Percentage))
	}
	defer C.free(unsafe.Pointer(src))
	defer C.free(unsafe.Pointer(save))
}

func MkDirs(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsExist(err) {
			return os.MkdirAll(path, os.ModePerm)
		}
	}
	return nil
}

func GetResizeFiles(srcPath, savePath, subPath string, curDepth, maxDepth int,
	rff ResizeFileFilter, fnc FileNameConversion) ([]ResizeFile, error) {
	if fnc == nil {
		return nil, fmt.Errorf("文件名转换不能为nil")
	}
	var filePaths []ResizeFile
	fi, err := os.Stat(srcPath)
	if err != nil {
		return nil, fmt.Errorf("文件信息读取失败 %s\t路径 %s\n", err, srcPath)
	}
	if fi.IsDir() {
		subFileInfos, err := ioutil.ReadDir(srcPath)
		if err != nil {
			return nil, fmt.Errorf("目录文件读取失败 %s\t路径 %s\n", err, srcPath)
		} else {
			for _, file := range subFileInfos {
				subSrcPath := fmt.Sprintf("%s\\%s", srcPath, file.Name())
				if file.IsDir() {
					if maxDepth > 0 && curDepth+1 >= maxDepth {
						continue
					}
					subFiles, err := GetResizeFiles(subSrcPath, savePath,
						subPath+file.Name(), curDepth+1, maxDepth, rff, fnc)
					if err == nil {
						filePaths = append(filePaths, subFiles...)
					} else {
						return nil, err
					}
				} else {
					if rff != nil && rff(file.Name()) {
						continue
					}
					saveFileName := fnc(file.Name())
					var saveDirPath string
					if savePath != "" {
						saveFileDirPath := fmt.Sprintf("%s\\%s", savePath, subPath)
						if err := MkDirs(saveFileDirPath); err != nil {
							return nil, err
						}
						saveDirPath = fmt.Sprintf("%s\\%s", saveFileDirPath, saveFileName)
					} else {
						saveDirPath = fmt.Sprintf("%s\\%s", srcPath, saveFileName)
					}
					filePaths = append(filePaths, ResizeFile{SrcPath: subSrcPath, SavePath: saveDirPath})
				}
			}
		}
	} else {
		saveFileName := fnc(fi.Name())
		srcDirPath := filepath.Dir(srcPath)
		var saveDirPath string
		if savePath != "" {
			saveFileDirPath := fmt.Sprintf("%s\\%s", savePath, subPath)
			if err := MkDirs(saveFileDirPath); err != nil {
				return nil, err
			}
			saveDirPath = fmt.Sprintf("%s\\%s", saveFileDirPath, saveFileName)
		} else {
			saveDirPath = fmt.Sprintf("%s\\%s", srcDirPath, saveFileName)
		}

		filePaths = append(filePaths, ResizeFile{SrcPath: srcPath, SavePath: saveDirPath})
	}
	return filePaths, nil
}
