package fileutil

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
)

func calcMD5(file string) string {

	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Println(err)
	}

	md5 := fmt.Sprintf("%x", h.Sum(nil))
	return md5
}
