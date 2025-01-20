/*
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

 Written by Frederic PONT.
 (c) Frederic Pont 2024
*/

package global

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// Declare a global variable
var Exclude = ReadExclude()

func ReadExclude() []string {
	// Open the CSV file
	file, err := os.Open("config/exclude.csv")
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
