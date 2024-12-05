package conf

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

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

// ReadRules read the "rules.csv" file with the regex rules to modifiy the file name output from the directory to scan name
func ReadRules() [][]string {
	// Open the CSV file
	file, err := os.Open("src/configuration/rules.csv")
	if err != nil {
		fmt.Println(err)
		return [][]string{}
	}
	// Close the file when the function returns
	defer file.Close()

	// Create a new csv.Reader
	reader := csv.NewReader(file)
	// Set the delimiter to TAB
	reader.Comma = '\t'
	// Set the comment character to '#'
	reader.Comment = '#'
	// Set the number of fields per record to -1, ie dir path in the first 2 columns
	reader.FieldsPerRecord = -1
	// ski header
	reader.Read()

	// Create an empty slice of []string
	var rules [][]string

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
			return [][]string{}
		}

		// Append the value to rules
		rules = append(rules, []string{line[0], line[1]})
	}

	return rules
}
