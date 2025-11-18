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

type Args struct {
	Algorithm string // file footprint algo
	NbCPU     int    // nb of threads
	NbLines   int    //nb of line writen to file before flush
}

func ParseDir(dir string, args Args) {
	// Setup output file
	outfile, err := os.Create("results/" + args.Algorithm + "_" + DatePrefix(FormatName(dir)+".tsv"))
	if err != nil {
		log.Println(err)
		return
	}
	defer outfile.Close()

	writer := csv.NewWriter(outfile)
	writer.Comma = '\t' // Tab delimiter

	// Channel of files to process
	filesCh := make(chan string, 8192)
	// Channel for hash results
	resultsCh := make(chan []string, 8192)

	var wg sync.WaitGroup

	// Start worker goroutines for hash calculation
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

				// Send result to writer
				resultsCh <- []string{signature, fileInfo.Name(), path}
			}
		}()
	}

	// Writer goroutine (flush every 1024 lines)
	writeWg := sync.WaitGroup{}
	writeWg.Add(1)
	go func() {
		defer writeWg.Done()
		lineCount := 0
		for data := range resultsCh {
			if err := writer.Write(data); err != nil {
				fmt.Printf("Error writing to CSV: %v\n", err)
			}
			lineCount++
			if lineCount >= args.NbLines {
				writer.Flush()
				lineCount = 0
			}
		}
		// Final flush to ensure all is written
		writer.Flush()
	}()

	// Walk directory and send files to channel
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
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
	close(resultsCh)
	writeWg.Wait()
}

func GetFileAndPath(fullPath string) (string, string) {
	dir := filepath.Dir(fullPath)
	fileName := filepath.Base(fullPath)
	return dir, fileName
}

func DatePrefix(name string) string {
	return time.Now().Format("2006-01-02_150405_") + name
}
