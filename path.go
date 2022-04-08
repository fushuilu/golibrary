package golibrary

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

const FilepathSeparator = string(filepath.Separator)

//
//  PathToCaller
//  @Description: 调用函数堆栈信息
//  @param skip 上溯的栈帧数
//  @return string
//
func PathToCaller(skip int) string {
	_, filename, _, _ := runtime.Caller(skip)
	return path.Dir(filename) + FilepathSeparator
}

//
//  PathGoLibrary
//  @Description: 获取当前 golibrary 所在的目录路径
//  @return string /.../github.com/fushuilu/golibrary/
//
func PathGoLibrary() string {
	return PathToCaller(0)
}

//
//  CallerDirectory
//  @Description: 调用此函数的文件所在的目录
//  @return string
//
func CallerDirectory() string {
	return PathToCaller(2)
}

//
//  PathSub
//  @Description: 目录截取
//  @param abPath 待截止目录的绝对
//  @param skip 忽略子目录的层次
//  @return string
//  @return error
//
func PathSub(abPath string, skip int) (string, error) {
	if abPath == "" {
		return "", errors.New("待截取目录不能为空")
	}
	if !strings.HasSuffix(abPath, FilepathSeparator) {
		abPath += FilepathSeparator
	}
	ns := strings.Split(abPath, FilepathSeparator)

	pl := len(ns) - skip - 1
	if pl < 1 {
		return "", errors.New("截取目录层次过长")
	}
	return strings.Join(ns[0:pl], FilepathSeparator) + FilepathSeparator, nil
}

// 判断文件是否存在
// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
func FileExist(path string) (bool, error) {
	if path == "" {
		return false, nil
	}
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false, err
	}
}
func IsFileExist(path string) bool {
	if exist, err := FileExist(path); err != nil {
		return false
	} else if exist {
		return true
	}
	return false
}

// 判断文件路径是否存在（只运行在 *unix 系统下），无法在 window 下运行
//func FileExistPosix(path string) bool {
//	err := syscall.Access(path, syscall.F_OK)
//	return !os.IsNotExist(err)
//}

const (
	Byte = 1
	KB   = Byte * 1024
	MB   = KB * 1024
	GB   = MB * 1024
)

func FilePath(pathDir, name string) string {

	pathEnd := strings.HasSuffix(pathDir, FilepathSeparator)
	nameStart := strings.HasPrefix(name, FilepathSeparator)

	if pathEnd {
		if nameStart {
			return pathDir + name[1:]
		} else {
			return pathDir + name
		}
	} else {
		if nameStart {
			return pathDir + name
		} else {
			return pathDir + FilepathSeparator + name
		}
	}
}

// 递归的创建目录
func CreateDirection(pathfile string, isDir bool) error {
	if strings.Index(pathfile, "..") > -1 {
		return errors.New("不允许创建相对目录")
	}
	dir := pathfile
	if !isDir {
		dir = filepath.Dir(pathfile)
	}
	if exist, err := FileExist(dir); err != nil {
		return err
	} else if !exist {
		root := PathGoLibrary()
		if len(pathfile) < len(root) {
			return errors.New("不允许创建超出 root 目录的目录")
		}
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}
