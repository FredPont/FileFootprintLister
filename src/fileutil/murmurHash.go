package fileutil

import (
	"fmt"
	"io"
	"os"

	"github.com/spaolacci/murmur3"
)

func calcMurmurHash64(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	hasher := murmur3.New64()
	if _, err := io.Copy(hasher, file); err != nil {
		panic(err)
	}
	murmurHash := fmt.Sprintf("%x", hasher.Sum64())
	return murmurHash
}
