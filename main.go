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
	fileutil "FileFootprintLister/src/fileutil"
	"FileFootprintLister/src/global"
	"flag"
	"fmt"
	"time"
)

func main() {
	fileutil.Title()

	t0 := time.Now()
	args := parseARG()
	fmt.Println("Parameters : ", args)

	allDirPath := conf.ReadAllPath()

	fmt.Println("Path containing the folowing strings will be excluded : ", global.Exclude)

	// start a new goroutine that runs the spinner function
	// Create a channel called stop
	stop := make(chan struct{})
	go fileutil.Spinner(stop) // enable spinner

	for i, dp := range allDirPath {
		t_start := time.Now()
		fmt.Println("\r", i+1, "/", len(allDirPath), " - ", dp, " is analysed...")
		// start current directory analysis
		fileutil.ParseDir(dp, args)
		if time.Since(t_start) < time.Second {
			time.Sleep(1 * time.Second) // sleep to enable file saving with date time prefix
		}
	}

	close(stop) // closing the channel stop the goroutine

	fmt.Println("\ndone !")
	fmt.Println("Elapsed time : ", time.Since(t0))
}

// parse arg of the command line and return the argument struct
func parseARG() fileutil.Args {
	args := fileutil.Args{}
	flag.StringVar(&args.Algorithm, "a", "md5", "algorithm to use. md5, xxhash, murmmur or sha256")
	flag.IntVar(&args.NbCPU, "n", 8, "number of CPUs for parallel file processing")
	flag.Parse()
	return args
}
