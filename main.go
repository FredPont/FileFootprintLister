/*
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

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
	fmt.Println("Parameters:", args)

	allDirPath := conf.ReadAllPath()
	fmt.Println("Excluded path patterns:", global.Exclude)

	// Spinner runs in its own goroutine; closing stop shuts it down.
	stop := make(chan struct{})
	go fileutil.Spinner(stop)

	for i, dp := range allDirPath {
		fmt.Printf("\r%d / %d  —  %s  is being analysed...\n", i+1, len(allDirPath), dp)
		fileutil.ParseDir(dp, args)
		// Output filenames now include a per-directory index suffix so they are
		// guaranteed to be unique even when two directories finish within the
		// same second. See DatePrefix in parsedir.go.
	}

	close(stop)

	fmt.Println("\ndone!")
	fmt.Println("Elapsed time:", time.Since(t0))
}

// parseARG parses CLI flags and returns an Args struct.
func parseARG() fileutil.Args {
	args := fileutil.Args{}
	flag.StringVar(&args.Algorithm, "a", "md5",
		"Hash algorithm: md5, xxhash, murmur, cityhash64, cityhash128, clickhouse64, clickhouse128, sha256")
	flag.IntVar(&args.NbCPU, "n", 8,
		"Number of parallel hash workers.\n"+
			"  HDD: keep at 1-2 to avoid seek overhead.\n"+
			"  SSD/NVMe: 4-16 gives best throughput.")
	flag.IntVar(&args.NbLines, "f", 1024,
		"Number of result rows buffered before flushing to disk")
	flag.Parse()
	return args
}
