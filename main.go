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

package main

import (
	conf "FileFootprintLister/src/configuration"
	"FileFootprintLister/src/fileutil"
	"flag"
	"fmt"
	"time"
)

type Args struct {
	Algorithm string
}

// ###########################################
func header() {
	fmt.Println("")
	fmt.Println("   ┌────────────────────────────────────────┐") // unicode U+250C
	fmt.Println("   │  FileFootprintLister (c)Frederic PONT  │")
	fmt.Println("   │    v20240511 - Free Software GNU GPL   │")
	fmt.Println("   └────────────────────────────────────────┘")
	//fmt.Println("")
}

// ###########################################
func main() {
	header()
	t0 := time.Now()
	//args := parseARG()

	allDirPath := conf.ReadAllPath()

	// start a new goroutine that runs the spinner function
	// Create a channel called stop
	stop := make(chan struct{})
	go fileutil.Spinner(stop) // enable spinner

	for _, dp := range allDirPath {
		t_start := time.Now()
		fmt.Println(dp, " is analysed...")
		fileutil.ParseDir(dp)
		if time.Since(t_start) < time.Second {
			time.Sleep(1 * time.Second) // sleep to enable file saving with date time prefix
		}
	}

	close(stop) // closing the channel stop the goroutine

	fmt.Println("\ndone !")
	fmt.Println("Elapsed time : ", time.Since(t0))
}

// parse arg of the command line and return the argument struct
func parseARG() Args {
	args := Args{}
	flag.StringVar(&args.Algorithm, "a", "md5", "algorithm to use")
	flag.Parse()

	return args
}
