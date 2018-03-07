package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	// get input from the CLI
	csvpath := flag.String("path", "./flights.csv", "File path")
	flag.Parse()

	content := ReadCSV(csvpath)

	headers := make([]string, 0)
	for _, head := range content[0] {
		headers = append(headers, head) // get the header values
	}

	content = content[1:] // slice the array in order to remove the header row as we already assigned it to the headers array.

	var buffer bytes.Buffer
	buffer = ConvertJSON(headers, content)

	path := GetPath() + "\\go-csvtojson" // temporary solution

	newFileName := filepath.Base(path)
	newFileName = newFileName[0:len(newFileName)-len(filepath.Ext(newFileName))] + ".json"
	r := filepath.Dir(path)
	filePath := filepath.Join(r, newFileName)

	SaveFile(&buffer, filePath)
}

// ReadCSV read csv file and prepare for parsing
func ReadCSV(csvpath *string) [][]string {
	f, err := os.Open(*csvpath) // automatic file upload will be added.
	check(err)
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1 // to avoid fieldcheckerror
	content, err := reader.ReadAll()
	check(err)

	return content
}

// ConvertJSON convert it to json
func ConvertJSON(headers []string, content [][]string) bytes.Buffer {
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

// SaveFile write output buffer to a file
func SaveFile(myFile *bytes.Buffer, path string) {
	if err := ioutil.WriteFile(path, myFile.Bytes(), os.FileMode(0644)); err != nil {
		panic(err)
	}
}

// GetPath return dir path
func GetPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func check(e error) {
	if e != nil {
		panic(e.Error())
	}
}
