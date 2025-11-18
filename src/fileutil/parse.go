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
	"FileFootprintLister/src/global"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	mutex sync.Mutex
)

// Args holds the argument of the software (md5, sha256)
type Args struct {
	Algorithm string
	NbCPU     int
}

func ParseDir(dir string, args Args) {
	// Setup output file
	outfile, err := os.Create("results/" + args.Algorithm + "_" + DatePrefix(FormatName(dir)+".tsv"))
	if err != nil {
		log.Println(err)
		return
	}
	defer outfile.Close()

	// Channel of files to process
	filesCh := make(chan string, 1024)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < args.NbCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range filesCh {
				fileInfo, err := os.Stat(path)
				if err != nil || fileInfo.IsDir() {
					continue
				}

				signature := ""
				switch args.Algorithm {
				case "sha256":
					signature = calcSHA256(path)
				case "xxhash":
					signature = calcXXHash64(path)
				case "murmur":
					signature = calcMurmurHash64(path)
				case "cityhash64":
					signature = calcCityHash64(path)
				case "cityhash128":
					signature = calcCityHash128(path)
				case "clickhouse64":
					signature = calcClickHouse64(path)
				case "clickhouse128":
					signature = calcClickHouse128(path)
				case "md5":
					signature = calcMD5(path)
				default:
					signature = calcMD5(path)
				}

				writeToCSV(outfile, []string{signature, fileInfo.Name(), path})
			}
		}()
	}

	// Walk directory and send files to channel
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Optionally handle exclude dirs!
			for _, ex := range global.Exclude {
				if strings.Contains(path, ex) {
					return filepath.SkipDir
				}
			}
			return nil
		}
		filesCh <- path
		return nil
	})

	if err != nil {
		fmt.Println("WalkDir Error:", err)
	}

	close(filesCh)
	wg.Wait()
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
