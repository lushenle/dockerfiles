package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func main() {
	dir := "./"
	renLogFileName := ".renlog"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Failed to read dir %s, the error is %v\n", dir, err)
		return
	}
	renameLog := ""
	ignoreFile := make(map[string]int)
	num := 0
L1:
	for _, f := range files {
		fileName := f.Name()
		if fileName == renLogFileName {
			continue
		}
		// 代表取文件名第一个字符
		if fileName[:1] == "." {
			fmt.Printf("Skip hidden file %s\n", fileName)
			continue
		}
		// 文件后缀名
		var imageExtensions = []string{".bmp"}
		// 只修改当前文件目录下文件
		if !f.IsDir() {
			dotIndex := strings.LastIndex(fileName, ".")
			//判断文件是一个合法的常规文件
			if dotIndex != -1 && dotIndex != 0 {
				extensionName := fileName[dotIndex:]
				ok, _ := inArray(extensionName, imageExtensions)
				if !ok {
					continue
				}
				fmt.Printf("开始处理文件: %s\n", fileName)
				if _, ok := ignoreFile[fileName]; ok {
					fmt.Printf("%s 文件不做改变，直接跳过\n", fileName)
					continue L1
				}
				// go里面没有do while循环体，只能用这种写法
				var newFileName string
			L2:
				for {
					// 注意：++ 和 --只能作为语句使用，不能作为表达式
					num++
					// int转string
					newFileName = strconv.Itoa(num)
					newFileName += extensionName
					if newFileName == fileName {
						ignoreFile[fileName] = 1
						continue L1
					}
					if !fileExists(newFileName) {
						break L2
					} else {
						ignoreFile[newFileName] = 1
						fmt.Printf("%s 已经被占用\n", newFileName)
					}
				}
				fmt.Printf("%s 文件即将被修改成 %s", fileName, newFileName)

				err = os.Rename(fileName, newFileName)
				fmt.Printf("，转化 %s 为 %s\n", fileName, newFileName)
				if err != nil {
					fmt.Printf("Failed to rename file %s to %s, the error is %v\n", fileName, newFileName, err)
					continue
				}
				renameLog += fmt.Sprintf("%s\t%s\n", fileName, newFileName)
			}
		}
	}
	fmt.Printf(renameLog)
	// ioutil.WriteFile(renLogFileName, []byte(renameLog), 0666)
}

func inArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

// 判断所给路径文件/文件夹是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func GetMD5Hash(text string) string {
	return GetByteMD5Hash([]byte(text))
}

func GetByteMD5Hash(content []byte) string {
	hasher := md5.New()
	hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil))
}
