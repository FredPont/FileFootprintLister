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

package fileutil

import (
	"encoding/csv"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Args struct {
	Algorithm string
}

func ParseDir(dir string, args Args) {
	// Create a file for writing
	outfile, err := os.Create("results/" + args.Algorithm + "_" + DatePrefix("output.tsv"))
	if err != nil {
		log.Println(err)
	}
	defer outfile.Close()

	// Create a CSV writer with a tab delimiter for TSV format
	writer := csv.NewWriter(outfile)
	writer.Comma = '\t' // Set the delimiter to tab

	err = filepath.WalkDir(dir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Process the file or directory here
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}

		startFootprintCalc(path, fileInfo, writer, args)

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
	// Flush writes any buffered data to the underlying io.Writer
	writer.Flush()

	// Check if there have been any errors during Write or Flush
	if err := writer.Error(); err != nil {
		panic(err) // Handle errors after flushing
	}
}

func writeLine(writer *csv.Writer, data []string) {
	// Write the []string as a row to the file
	err := writer.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func startFootprintCalc(path string, fileInfo fs.FileInfo, writer *csv.Writer, args Args) {
	if !fileInfo.IsDir() {
		var signature string
		// compute file footprint
		if args.Algorithm == "sha256" {
			signature = calcSHA256(path)
		} else {
			signature = calcMD5(path)
		}

		//fmt.Println("Visited:", path, " ", signature)
		writeLine(writer, []string{signature, path})
	}
}

func GetFileAndPath(fullPath string) (string, string) {
	dir := filepath.Dir(fullPath)
	fileName := filepath.Base(fullPath)
	return dir, fileName
}

// DatePrefix, prefix a string with current date and time
func DatePrefix(name string) string {
	return time.Now().Format("2006-01-02_150405_") + name
}
