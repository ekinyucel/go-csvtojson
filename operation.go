package main

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func trackFiles() {
	files, err := ioutil.ReadDir(folderName)
	if err != nil {
		logger.Printf("Error: %s", err.Error())
	}

	for _, file := range files {
		if getInputFileFormat(file, fileType) && !isFileProcessed(&fileList, file.Name()) {
			fileList = append(fileList, File{filename: file.Name(), processed: false})
		}
	}

	fileChannel <- fileList
}

func processFile(file *File) {
	startTime := time.Now()
	filename := file.filename

	if fileType == CSV {
		processCSV(filename)
	}

	file.processed = true
	endTime := time.Now()
	logger.Println(filename, " processed in ", endTime.Sub(startTime))
}

func processCSV(filename string) {
	content := readCSV(&filename)

	headers := make([]string, len(content[0]))
	for i, head := range content[0] {
		headers[i] = head // get the header values
	}
	content = content[1:] // slice the array in order to remove the header row as we already assigned it to the headers array.

	var buffer bytes.Buffer
	buffer = convertJSON(headers, content)

	logger.Println("is json ", isJSON(buffer.Bytes()))

	newFileName := filename + strconv.FormatInt(time.Now().Unix(), 10)
	newFileName = newFileName[0:len(newFileName)-len(filepath.Ext(newFileName))] + "." + targetType
	r := filepath.Dir(folderName)
	filePath := filepath.Join(r, newFileName)

	if err := saveFile(&buffer, filePath); err != nil {
		logger.Printf("error: %v", err)
	}
}

func readCSV(csvpath *string) [][]string {
	f, err := os.Open(folderName + *csvpath) // automatic file upload will be added.
	if err != nil {
		logger.Printf("Error: %s", err.Error())
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1 // to avoid fieldcheckerror
	content, err := reader.ReadAll()
	if err != nil {
		logger.Printf("Error: %s", err.Error())
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
				if headers[x] != "   " { // if the header is empty, then discard the entire column
					buffer.WriteString(`"` + headers[x] + `":`)
					_, err := strconv.ParseFloat(y, 32)
					_, err2 := strconv.ParseBool(y)
					if err == nil {
						buffer.WriteString((`"` + strings.TrimSpace(y) + `"`))
					} else if err2 == nil {
						buffer.WriteString((`"` + strings.TrimSpace(strings.ToLower(y)) + `"`))
					} else {
						buffer.WriteString((`"` + strings.Replace(y, "\n", "", -1) + `"`))
					}

					if x < len(d)-2 { // I wrote len(d)-2 in order to avoid extra comma after the last field. it had an issue with extra comma in the end
						buffer.WriteString(string(","))
					}
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
