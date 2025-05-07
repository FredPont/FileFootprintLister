package fileutil

import (
	"fmt"
	"io"
	"os"

	"github.com/cespare/xxhash"
)

func calcXXHash64(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	hasher := xxhash.New()
	if _, err := io.Copy(hasher, file); err != nil {
		panic(err)
	}

	XXHash64 := fmt.Sprintf("%x", hasher.Sum64())
	return XXHash64
}
