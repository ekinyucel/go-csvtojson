package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var logger = log.New(os.Stdout, "converter: ", log.LstdFlags)
var fileList []string
var wg sync.WaitGroup

func main() {
	startTime := time.Now()
	fileType := "csv"

	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if getInputFileFormat(file, fileType) {
			fileList = append(fileList, file.Name())
		}
	}

	if len(fileList) == 0 {
		logger.Printf("no %s file is present", fileType)
		return
	}

	done := make(chan bool)

	wg.Add(len(fileList))

	for i := range fileList {
		go processFile(&fileList[i], done) // if there are more than one file for this type file format, then it speeds up the process.
	}

	wg.Wait() // waiting here until all goroutines are finished their execution

	endTime := time.Now()
	fmt.Println(endTime.Sub(startTime))
}

// GetPath return dir path
func GetPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
