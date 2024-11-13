package conf

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Conf struct {
	NbCPU int `json:"nbCPU"`
}

// ReadConf read json conf file
func ReadConf() Conf {
	fname := "config/conf.json"
	var cs Conf
	fp, err := os.Open(fname)
	if err != nil {
		fmt.Println(err)
	}
	defer fp.Close()
	bytes, err := io.ReadAll(fp)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(bytes, &cs)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(cs)
	return cs
}

func ReadAllPath() []string {
	// Open the CSV file
	file, err := os.Open("config/path.csv")
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	// Close the file when the function returns
	defer file.Close()

	// Create a new csv.Reader
	reader := csv.NewReader(file)
	// Set the delimiter to TAB
	//reader.Comma = '\t'
	// Set the comment character to '#'
	reader.Comment = '#'
	// Set the number of fields per record to 1, ie dir path in the first column
	reader.FieldsPerRecord = 1

	// Create an empty slice of []string
	var allPath []string

	// Loop through the remaining lines
	for {
		// Read a line
		line, err := reader.Read()
		// Check the error value
		if err != nil {
			// Break the loop when the end of the file is reached
			if err == io.EOF {
				break
			}
			// Print the error otherwise
			fmt.Println(err)
			return []string{}
		}

		// Append the value to allPath
		allPath = append(allPath, line[0])
	}

	return allPath
}
