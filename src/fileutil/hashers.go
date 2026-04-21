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
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"log"
	"os"

	"github.com/cespare/xxhash"
	"github.com/go-faster/city"
	"github.com/spaolacci/murmur3"
)

// ── Streaming hashers ─────────────────────────────────────────────────────────
// These wrap stdlib / third-party types so they satisfy hasherIface and can be
// used with the shared bufPool via hashStream() in parsedir.go.

// md5Hasher wraps crypto/md5 — already satisfies hash.Hash / hasherIface.
type md5Hasher struct{ hash.Hash }

func newMD5Hasher() hasherIface { return &md5Hasher{md5.New()} }

// sha256Hasher wraps crypto/sha256.
type sha256Hasher struct{ hash.Hash }

func newSHA256Hasher() hasherIface { return &sha256Hasher{sha256.New()} }

// xxHasher wraps cespare/xxhash via the hash.Hash64 interface it returns.
// xxhash.New() returns a hash.Hash (which embeds io.Writer and Sum), so we
// store it as hash.Hash and override Sum to emit the hex-encoded 64-bit value.
type xxHasher struct {
	h   hash.Hash   // underlying xxhash writer (also implements hash.Hash64)
	h64 interface { // we need Sum64(); cast once at construction time
		Sum64() uint64
	}
}

func newXXHasher() hasherIface {
	d := xxhash.New() // *xxhash.Digest implements both hash.Hash and Sum64()
	return &xxHasher{h: d, h64: d}
}

func (x *xxHasher) Write(p []byte) (int, error) { return x.h.Write(p) }

func (x *xxHasher) Sum(b []byte) []byte {
	return []byte(fmt.Sprintf("%x", x.h64.Sum64()))
}

// murmurHasher wraps spaolacci/murmur3 via the hash.Hash32/64 interface.
// murmur3.New64() returns a hash.Hash64 — we store the Sum64 capability
// separately via a minimal interface to avoid embedding the unexported type.
type murmurHasher struct {
	h   hash.Hash
	h64 interface {
		Sum64() uint64
	}
}

func newMurmurHasher() hasherIface {
	d := murmur3.New64() // implements hash.Hash64
	return &murmurHasher{h: d, h64: d}
}

func (m *murmurHasher) Write(p []byte) (int, error) { return m.h.Write(p) }

func (m *murmurHasher) Sum(b []byte) []byte {
	return []byte(fmt.Sprintf("%x", m.h64.Sum64()))
}

// ── Non-streaming hashers (CityHash / ClickHouse) ────────────────────────────
// The go-faster/city library does not expose an io.Writer interface, so we
// must read the whole file into memory. A bufPool buffer is used to avoid
// extra allocations; for files larger than 1 MiB the buffer is grown once.

func calcCityHash64(filePath string) string {
	data := readFileBytes(filePath)
	if data == nil {
		return ""
	}
	return fmt.Sprintf("%x", city.Hash64(data))
}

func calcCityHash128(filePath string) string {
	data := readFileBytes(filePath)
	if data == nil {
		return ""
	}
	return fmt.Sprintf("%x", city.Hash128(data))
}

func calcClickHouse64(filePath string) string {
	data := readFileBytes(filePath)
	if data == nil {
		return ""
	}
	return fmt.Sprintf("%x", city.CH64(data))
}

func calcClickHouse128(filePath string) string {
	data := readFileBytes(filePath)
	if data == nil {
		return ""
	}
	return fmt.Sprintf("%x", city.CH128(data))
}

// readFileBytes reads an entire file efficiently.
// For files that fit in the pool buffer (≤ 1 MiB) the pooled slice is used;
// larger files fall back to os.ReadFile which allocates exactly once.
func readFileBytes(filePath string) []byte {
	info, err := os.Stat(filePath)
	if err != nil {
		log.Println("stat:", err)
		return nil
	}

	if info.Size() <= hashChunkSize {
		// Small file: read into pooled buffer to avoid allocation
		bufPtr := bufPool.Get().(*[]byte)
		defer bufPool.Put(bufPtr)
		buf := (*bufPtr)[:info.Size()]

		f, err := os.Open(filePath)
		if err != nil {
			log.Println("open:", err)
			return nil
		}
		defer f.Close()

		if _, err := io.ReadFull(f, buf); err != nil {
			log.Println("read:", err)
			return nil
		}
		// Return a copy — the pool buffer is returned to the pool above
		out := make([]byte, len(buf))
		copy(out, buf)
		return out
	}

	// Large file: single allocation via os.ReadFile
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("ReadFile:", err)
		return nil
	}
	return data
}
