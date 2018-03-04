package main

import (
	"encoding/csv"
	"os"
)

func main() {
	f, err := os.Open("./flights2.csv")
	check(err)
	defer f.Close()

	reader := csv.NewReader(f)
	//reader.FieldsPerRecord = -1

	content, err := reader.ReadAll()

	check(err)

	headers := make([]string, 0)
	for _, head := range content[0] {
		headers = append(headers, head) // get the header values
	}

	content = content[1:]
	/*for {
		record, err := r.
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record)
	}*/
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
