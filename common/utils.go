package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func IsFileExist(file string) bool {
	fi, err := os.Stat(file)
	if err != nil {
		return os.IsExist(err)
	}
	return !fi.IsDir()
}

func ReadLines(file string) ([]string, error) {
	if !IsFileExist(file) {
		return nil, fmt.Errorf("文件不存在: %s", file)
	}
	f, err := os.OpenFile(file, os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("文件打开失败, %s", err)
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("文件读取失败 %s", err)
	}
	return strings.Split(string(bytes), "\n"), nil
}
