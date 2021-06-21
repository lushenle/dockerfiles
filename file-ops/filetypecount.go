package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	fptCount            int
	normalFptCount      int
	exceptionalFptCount int
	bmpCount            int
	pngCount            int
	jpgCount            int
	rarCount            int
	otherCount          int
)

func count(path string) {
	err := filepath.Walk(path, func(path string, file os.FileInfo, err error) error {
		if file.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(file.Name()), ".fptx") || strings.HasSuffix(strings.ToLower(file.Name()), ".fpt") ||
			strings.HasSuffix(strings.ToLower(file.Name()), ".fpt.zip") || strings.HasSuffix(strings.ToLower(file.Name()), ".fptx.zip") {
			fptCount++
			if file.Size() >= 80000 {
				normalFptCount++
			} else {
				exceptionalFptCount++
			}
		} else if strings.HasSuffix(strings.ToLower(file.Name()), ".bmp") {
			bmpCount++
		} else if strings.HasSuffix(strings.ToLower(file.Name()), ".png") {
			pngCount++
		} else if strings.HasSuffix(strings.ToLower(file.Name()), ".jpg") {
			jpgCount++
		} else if strings.HasSuffix(strings.ToLower(file.Name()), ".rar") {
			rarCount++
		} else {
			otherCount++
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("FPTFiles: %d\nNormalFptCount: %d\nExceptionalFptCount: %d\nBMPFiles: %d\nPNGFiles: %d\nJPGFiles: %d\nRARFiles: %d\nOtherFiles: %d\n",
		fptCount, normalFptCount, exceptionalFptCount, bmpCount, pngCount, jpgCount, rarCount, otherCount)
}

func main() {
	count(os.Args[1])
}
