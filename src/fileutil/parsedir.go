/*
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Written by Frederic PONT.
(c) Frederic Pont 2024
*/

package fileutil

import (
	"FileFootprintLister/src/global"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// hashChunkSize is the read-buffer size used by all streaming hashers.
// 1 MiB matches rclone's default and gives a measurable throughput gain
// over the 32 KiB default used by io.Copy.
const hashChunkSize = 1 << 20 // 1 MiB

// bufPool is a process-wide pool of reusable 1 MiB byte slices.
// Reusing buffers avoids repeated heap allocations in the hash workers,
// which matters when hashing thousands of small files.
var bufPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, hashChunkSize)
		return &b
	},
}

// Args holds the CLI parameters for a scan run.
type Args struct {
	Algorithm string // hash algorithm name
	NbCPU     int    // number of hash worker goroutines
	NbLines   int    // results are flushed to disk every NbLines rows
}

// ParseDir walks dir, hashes every file with the chosen algorithm using
// NbCPU workers, and writes a TSV result file in results/.
func ParseDir(dir string, args Args) {
	outfile, err := os.Create("results/" + args.Algorithm + "_" + DatePrefix(FormatName(dir)+".tsv"))
	if err != nil {
		log.Println(err)
		return
	}
	defer outfile.Close()

	writer := csv.NewWriter(outfile)
	writer.Comma = '\t'

	// filesCh carries file paths from the walker to the hash workers.
	// A generous buffer keeps the walker from blocking while workers are busy.
	filesCh := make(chan string, args.NbCPU*64)

	// resultsCh carries completed [signature, name, path] rows to the writer.
	resultsCh := make(chan []string, args.NbCPU*64)

	// ── Hash workers ─────────────────────────────────────────────────────────
	var wg sync.WaitGroup
	for i := 0; i < args.NbCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range filesCh {
				info, err := os.Stat(path)
				if err != nil || info.IsDir() {
					continue
				}
				sig := computeHash(path, args.Algorithm)
				if sig == "" {
					continue
				}
				resultsCh <- []string{sig, info.Name(), path}
			}
		}()
	}

	// ── Writer goroutine ──────────────────────────────────────────────────────
	var writeWg sync.WaitGroup
	writeWg.Add(1)
	go func() {
		defer writeWg.Done()
		lineCount := 0
		for row := range resultsCh {
			if err := writer.Write(row); err != nil {
				fmt.Printf("CSV write error: %v\n", err)
			}
			lineCount++
			if lineCount >= args.NbLines {
				writer.Flush()
				lineCount = 0
			}
		}
		writer.Flush()
	}()

	// ── Directory walker (single goroutine — disk seeks are sequential) ───────
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
		fmt.Println("WalkDir error:", err)
	}

	close(filesCh)
	wg.Wait()
	close(resultsCh)
	writeWg.Wait()
}

// computeHash dispatches to the appropriate algorithm and returns the hex
// digest, or "" on error. Streaming hashers use the shared bufPool.
func computeHash(path, algorithm string) string {
	switch algorithm {
	case "sha256":
		return hashStream(path, newSHA256)
	case "xxhash":
		return hashStream(path, newXXHash64)
	case "murmur":
		return hashStream(path, newMurmur64)
	case "cityhash64":
		return calcCityHash64(path) // non-streaming — reads whole file
	case "cityhash128":
		return calcCityHash128(path)
	case "clickhouse64":
		return calcClickHouse64(path)
	case "clickhouse128":
		return calcClickHouse128(path)
	case "md5":
		return hashStream(path, newMD5)
	default:
		return hashStream(path, newMD5)
	}
}

// hashStream is a generic streaming hasher. It opens path, borrows a buffer
// from bufPool to read in 1 MiB chunks, and returns the hex digest.
// The constructor fn returns a fresh io.Writer that also implements Sum() via
// the hasherIface interface below.
func hashStream(path string, constructor func() hasherIface) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open:", err)
		return ""
	}
	defer f.Close()

	h := constructor()

	// Borrow a buffer from the pool; return it when done.
	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)
	buf := *bufPtr

	if _, err := io.CopyBuffer(h, f, buf); err != nil {
		log.Println("hash read:", err)
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// hasherIface is satisfied by any hash that supports streaming writes and
// can return its digest via Sum(nil). All standard library hashes (md5,
// sha256) and xxhash/murmur3 implement this interface.
type hasherIface interface {
	io.Writer
	Sum(b []byte) []byte
}

// ── Constructor wrappers ──────────────────────────────────────────────────────
// Each returns a hasherIface so hashStream can be generic.

func newMD5() hasherIface {
	// crypto/md5 is imported in md5.go — reuse the same import
	return newMD5Hasher()
}

func newSHA256() hasherIface {
	return newSHA256Hasher()
}

func newXXHash64() hasherIface {
	return newXXHasher()
}

func newMurmur64() hasherIface {
	return newMurmurHasher()
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func GetFileAndPath(fullPath string) (string, string) {
	return filepath.Dir(fullPath), filepath.Base(fullPath)
}

func DatePrefix(name string) string {
	return time.Now().Format("2006-01-02_150405_") + name
}
