package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// File struct contains
type File struct {
	filename  string
	processed bool
}

// todo find a better way to search through the slice
func isFileProcessed(list *[]File, filename string) bool {
	for _, f := range *list {
		if f.filename == filename {
			return f.processed
		}
	}
	return false
}

func getInputFileFormat(fileName os.FileInfo, formatType string) bool {
	if formatType == "" {
		return false
	}
	name := fileName.Name()
	fileFormat := strings.Join(strings.Split(name, ".")[1:], ".")
	if fileFormat == formatType {
		return true
	}
	return false
}

func saveFile(myFile *bytes.Buffer, path string) {
	if err := ioutil.WriteFile(path, myFile.Bytes(), os.FileMode(0644)); err != nil {
		logger.Printf("Error: %s", err.Error())
	}
}

func getPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Printf("Error: %s", err.Error())
	}
	return dir
}
