package main

import (
	"os"
	"log"
	"path/filepath"
	"strings"
)

func rename(path string) {
	err := filepath.Walk(path, func(path string, file os.FileInfo, err error) error {
		if file.IsDir() {
			return nil
		}
		if strings.HasPrefix(file.Name(), "segment_")  {
			// rename
			os.Rename(file.Name(), strings.Replace(file.Name(),"segment_","",1))
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main()  {
	rename(os.Args[1])
}