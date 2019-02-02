package main

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var splitDelimiter = "."

// File struct contains file operations
type File struct {
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
