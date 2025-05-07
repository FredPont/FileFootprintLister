package fileutil

import (
	"fmt"
	"log"
	"os"

	"github.com/go-faster/city"
)

func calcCityHash64(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Compute CityHash
	hash64 := city.Hash64(data)

	cityHash := fmt.Sprintf("%x", hash64)
	return cityHash
}

func calcCityHash128(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Compute CityHash
	hash128 := city.Hash128(data)

	cityHash := fmt.Sprintf("%x", hash128)
	return cityHash
}

func calcClickHouse64(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Compute CityHash
	hash64 := city.CH64(data)

	cityHash := fmt.Sprintf("%x", hash64)
	return cityHash
}

func calcClickHouse128(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Compute CityHash
	hash128 := city.CH128(data)

	cityHash := fmt.Sprintf("%x", hash128)
	return cityHash
}
