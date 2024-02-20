package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/robotomize/go-bf/bf/bfbits"
)

func main() {
	filter := bfbits.NewBloomFilter(100_000_000, bfbits.WithHashNum(3))

	buf := make([]byte, 8)
	b := make([][]byte, 5)

	for i := 0; i < 70_000_000; i++ {
		read, err := rand.Read(buf)
		if err != nil {
			return
		}
		key := []byte(hex.EncodeToString(buf[:read]))
		if err := filter.Add(key); err != nil {
			log.Fatal(err)
		}
		if i < len(b) {
			b[i] = key
		}
	}

	for i := 0; i < len(b); i++ {
		contains, err := filter.Contains(b[i])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(contains)
	}
}
