package main

import (
	"fmt"

	"github.com/robotomize/go-bf/bf/bfbits"
)

func main() {
	bf := bfbits.NewBloomFilter(10000, bfbits.WithHashNum(1))

	if err := bf.Add([]byte("apple")); err != nil {
		fmt.Printf("Add() error = %v", err)
	}

	items := []string{"banana", "cherry", "date"}
	for _, item := range items {
		if err := bf.Add([]byte(item)); err != nil {
			fmt.Printf("Add() error = %v", err)
		}
	}
}
