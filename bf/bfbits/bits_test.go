package bfbits

import (
	"crypto/rand"
	"encoding/hex"
	"runtime"
	"sync"
	"testing"
)

func BenchmarkBloomFilter_Contains(b *testing.B) {
	filter := NewBloomFilter(100_000_000, WithHashNum(3))
	wCh := make(chan []byte, 1)
	rCh := make(chan []byte, 1)
	wChGroup := sync.WaitGroup{}
	wChGroup.Add(runtime.NumCPU())

	buf := make([]byte, 16)
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			defer wChGroup.Done()
			for i := 0; i < 10_000_000/runtime.NumCPU(); i++ {
				read, err := rand.Read(buf)
				if err != nil {
					b.Error(err)
					return
				}

				chunk := make([]byte, hex.EncodedLen(16))
				n := hex.Encode(chunk, buf[:read])
				wCh <- chunk[:n]
				rCh <- chunk[:n]
			}
		}()
	}

	go func() {
		wChGroup.Wait()
		close(wCh)
		close(rCh)
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for bytes := range wCh {
			if err := filter.Add(bytes); err != nil {
				b.Error(err)
				return
			}
		}
	}()

	go func() {
		wg.Done()
		for bytes := range rCh {
			if _, err := filter.Contains(bytes); err != nil {
				b.Error(err)
				return
			}
		}
	}()

	wg.Wait()
}

func TestBloomFilter_Add(t *testing.T) {
	bf := NewBloomFilter(10000, WithHashNum(1))

	if err := bf.Add([]byte("apple")); err != nil {
		t.Errorf("Add() error = %v", err)
	}

	items := []string{"banana", "cherry", "date"}
	for _, item := range items {
		if err := bf.Add([]byte(item)); err != nil {
			t.Errorf("Add() error = %v", err)
		}
	}
}

func TestBloomFilter_Contains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		size     int
		hashNum  int
		items    []string
		checked  string
		expected bool
	}{
		{
			name:     "test_single_true",
			size:     10000,
			hashNum:  1,
			items:    []string{"apple"},
			checked:  "apple",
			expected: true,
		},
		{
			name:     "test_multiple_all_true",
			size:     10000,
			hashNum:  1,
			items:    []string{"apple", "banana", "cherry"},
			checked:  "banana",
			expected: true,
		},
		{
			name:     "test_multiple_with_false",
			size:     10000,
			hashNum:  1,
			items:    []string{"apple", "banana", "date"},
			checked:  "guava",
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				bf := NewBloomFilter(tt.size, WithHashNum(tt.hashNum))

				for _, item := range tt.items {
					if err := bf.Add([]byte(item)); err != nil {
						t.Errorf("Add() error = %v", err)
					}
				}

				contains, err := bf.Contains([]byte(tt.checked))
				if err != nil {
					t.Errorf("Contains() error = %v", err)
				}
				if contains != tt.expected {
					t.Errorf("Contains() got = %v, want %v", contains, tt.expected)
				}
			},
		)
	}
}
