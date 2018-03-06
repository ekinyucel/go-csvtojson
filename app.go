package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("./flights2.csv")
	check(err)
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1 // to avoid fieldcheckerror
	content, err := reader.ReadAll()

	check(err)

	headers := make([]string, 0)
	for _, head := range content[0] {
		headers = append(headers, head) // get the header values
	}

	content = content[1:]

	var buffer bytes.Buffer
	buffer.WriteString(string("["))

	for i, d := range content {
		buffer.WriteString(string("{"))

		for x, y := range d {
			if x < len(headers)-1 {
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

				if x < len(d)-1 {
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
	//raw := json.RawMessage(buffer.String())
	output, _ := json.MarshalIndent(buffer, "", " ")
	//fmt.Println(raw)
	fmt.Println(&buffer)
	fmt.Println(output)

	path := "C:/Users/eyucel/go/src/go-csvtojson/"

	newFileName := filepath.Base(path)
	newFileName = newFileName[0:len(newFileName)-len(filepath.Ext(newFileName))] + ".json"
	r := filepath.Dir(path)
	filePath := filepath.Join(r, newFileName)

	SaveFile(&buffer, filePath)
}

func SaveFile(myFile *bytes.Buffer, path string) {
	if err := ioutil.WriteFile(path, myFile.Bytes(), os.FileMode(0644)); err != nil {
		panic(err)
	}
}

/*func GetCurrentPath() *string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}*/

func check(e error) {
	if e != nil {
		panic(e.Error())
	}
}
