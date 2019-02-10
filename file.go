package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var splitDelimiter = "."

// File struct contains
type File struct {
	filename  string
	processed bool
}

func trackFiles(c chan []File) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if getInputFileFormat(file, fileType) && !isFileProcessed(&fileList, file.Name()) {
			fileList = append(fileList, File{filename: file.Name(), processed: false})
		}
	}
	c <- fileList
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

func processFile(file *File) {
	startTime := time.Now()
	filename := file.filename

	//defer wg.Done()
	content := readCSV(&filename)

	headers := make([]string, 0)
	for _, head := range content[0] {
		headers = append(headers, head) // get the header values
	}
	content = content[1:] // slice the array in order to remove the header row as we already assigned it to the headers array.

	var buffer bytes.Buffer
	buffer = convertJSON(headers, content)

	path := getPath() + "\\go-csvtojson" // temporary solution

	newFileName := filename + strconv.FormatInt(time.Now().Unix(), 10)
	newFileName = newFileName[0:len(newFileName)-len(filepath.Ext(newFileName))] + ".json"
	r := filepath.Dir(path)
	filePath := filepath.Join(r, newFileName)

	saveFile(&buffer, filePath)

	file.processed = true
	endTime := time.Now()
	fmt.Println(endTime.Sub(startTime))
}

func getInputFileFormat(fileName os.FileInfo, formatType string) bool {
	if formatType == "" {
		return false
	}
	name := fileName.Name()
	fileFormat := strings.Join(strings.Split(name, splitDelimiter)[1:], splitDelimiter)
	if fileFormat == formatType {
		return true
	}
	return false
}

func readCSV(csvpath *string) [][]string {
	f, err := os.Open("./" + *csvpath) // automatic file upload will be added.
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1 // to avoid fieldcheckerror
	content, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	return content
}

func convertJSON(headers []string, content [][]string) bytes.Buffer {
	var buffer bytes.Buffer

	buffer.WriteString(string("["))

	for i, d := range content {
		buffer.WriteString(string("{"))

		for x, y := range d {
			if x < len(headers)-1 { // check if we are in the limits of headers array when the iteration happens.
				buffer.WriteString(`"` + headers[x] + `":`)
				_, err := strconv.ParseFloat(y, 32)
				_, err2 := strconv.ParseBool(y)
				if err == nil {
					buffer.WriteString(y)
				} else if err2 == nil {
					buffer.WriteString(strings.ToLower(y))
				} else {
					buffer.WriteString((`"` + y + `"`))
				}

				if x < len(d)-4 { // I wrote len(d)-4 in order to avoid extra comma after the last field. it had an issue with extra comma at the end
					buffer.WriteString(string(","))
				}
			}
		}

		buffer.WriteString(string("}"))
		if i < len(content)-1 {
			buffer.WriteString(",")
		}
	}

	buffer.WriteString(string("]"))
	return buffer
}

func saveFile(myFile *bytes.Buffer, path string) {
	if err := ioutil.WriteFile(path, myFile.Bytes(), os.FileMode(0644)); err != nil {
		panic(err)
	}
}

func getPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}