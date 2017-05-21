package ucfile

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	DEFAULT_COMPRESS_TYPE = "gz"
)

func IsFile(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func IsDir(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// 获取绝对路径
func AbsolutePath(path string) string {
	pathAbs, err := filepath.Abs(path)
	if err != nil {
		pathAbs = path
		fmt.Println(err)
	}
	return pathAbs
}

// 遍历目录下的所有文件，不进入下一级目录搜索，可以匹配前缀过滤
func ScanDir(dirPath string, prefix string, handler func(string)) error {
	prefix = strings.ToUpper(prefix)
	pathSeparator := string(os.PathSeparator)
	fileInfos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, info := range fileInfos {
		if info.IsDir() {
			continue
		}
		if strings.HasPrefix(strings.ToUpper(info.Name()), prefix) {
			filename := dirPath + pathSeparator + info.Name()
			if handler != nil {
				handler(filename)
			} else {
				fmt.Println(filename)
			}
		}
	}
	return nil
}

// 递归遍历指定目录及所有子目录下的所有文件，可以匹配前缀过滤
func ScanDirRecursive(dirPath, prefix string, handler func(string)) error {
	prefix = strings.ToUpper(prefix)
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasPrefix(strings.ToUpper(info.Name()), prefix) {
			if handler != nil {
				handler(path)
			} else {
				fmt.Println(path)
			}
		}
		return nil
	})
	return err
}

// 读取文件中的数据，一次读取一行
func ScanFile(filePath string, handler func(int64, int64, string)) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	var lineNum int64 = 1
	var offset int64 = 0
	var scanner *bufio.Scanner
	// 支持普通文件和.gz压缩文件
	if strings.HasSuffix(filePath, DEFAULT_COMPRESS_TYPE) {
		gzipReader, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		defer gzipReader.Close()
		scanner = bufio.NewScanner(gzipReader)
	} else {
		scanner = bufio.NewScanner(f)
	}
	for scanner.Scan() {
		line := scanner.Text()
		if handler != nil {
			handler(lineNum, offset, line)
		} else {
			fmt.Println(line)
		}
		lineNum++
		offset += int64(len(line))
	}
	return nil
}
