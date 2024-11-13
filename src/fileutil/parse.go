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
	conf "FileFootprintLister/src/configuration"
	"encoding/csv"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	mutex sync.Mutex
)

type Args struct {
	Algorithm string
}

func ParseDir(dir string, args Args) {
	var wg sync.WaitGroup
	maxGoroutines := conf.Config.NbCPU
	// create a channel with the max number of job allowed
	ch := make(chan struct{}, maxGoroutines)

	// Create a file for writing
	outfile, err := os.Create("results/" + args.Algorithm + "_" + DatePrefix("output.tsv"))
	if err != nil {
		log.Println(err)
	}
	defer outfile.Close()

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

		wg.Add(1) // Increment the WaitGroup counter
		ch <- struct{}{}
		go ParallelFootprintCalc(path, fileInfo, outfile, args, &wg, ch)

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}

	wg.Wait() // Wait for all goroutines to finish

}

func ParallelFootprintCalc(path string, fileInfo fs.FileInfo, outfile *os.File, args Args, wg *sync.WaitGroup, ch chan struct{}) {
	defer func() {
		<-ch // get the token back to free up a slot
		wg.Done()
	}()
	if !fileInfo.IsDir() {
		var signature string
		// compute file footprint
		switch args.Algorithm {
		case "sha256":
			signature = calcSHA256(path)
		default:
			signature = calcMD5(path)
		}

		writeToCSV(outfile, []string{signature, fileInfo.Name(), path})
	}
}

// writeToCSV write data to file concurrently
func writeToCSV(file *os.File, data []string) {

	// Lock the mutex to ensure exclusive access to the file
	mutex.Lock()
	defer mutex.Unlock()

	writer := csv.NewWriter(file)
	writer.Comma = '\t' // Set the delimiter to tab
	defer writer.Flush()

	// Write the data to the CSV file
	if err := writer.Write(data); err != nil {
		fmt.Printf("Error writing to CSV: %v\n", err)
		return
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
