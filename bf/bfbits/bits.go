package bfbits

import (
	"fmt"
	"hash"
	"hash/fnv"
)

type BloomFilter interface {
	Add(item []byte) error
	Contains(item []byte) (bool, error)
}

type HashFunc func() hash.Hash64

var defaultHashFunc = func() hash.Hash64 {
	return fnv.New64a()
}

var defaultHashNum = 1

type Option func(filter *bloomFilter)

func WithHashNum(n int) Option {
	return func(filter *bloomFilter) {
		filter.hashNum = n
	}
}

func WithHasher(f HashFunc) Option {
	return func(filter *bloomFilter) {
		filter.hasher = f
	}
}

func NewBloomFilter(size int, opts ...Option) BloomFilter {
	bf := bloomFilter{
		size:    size,
		bits:    make([]uint64, (size+63)/64),
		hasher:  defaultHashFunc,
		hashNum: defaultHashNum,
	}

	for _, o := range opts {
		o(&bf)
	}

	return &bf
}

type bloomFilter struct {
	size    int
	hashNum int
	hasher  func() hash.Hash64
	bits    []uint64
}

func (bf *bloomFilter) Add(item []byte) error {
	h := bf.hasher()
	for i := 0; i < bf.hashNum; i++ {
		if _, err := h.Write(item); err != nil {
			return fmt.Errorf("hasher Write: %w", err)
		}

		idx := int(h.Sum64() % uint64(bf.size))
		bf.bits[idx/64] |= 1 << (idx % 64)
	}

	return nil
}

func (bf *bloomFilter) Contains(item []byte) (bool, error) {
	h := bf.hasher()
	if _, err := h.Write(item); err != nil {
		return false, fmt.Errorf("hasher Write: %w", err)
	}

	idx := int(h.Sum64() % uint64(bf.size))

	r := bf.bits[idx/64]&(1<<(idx%64)) != 0
	for i := 1; i < bf.hashNum; i++ {
		if _, err := h.Write(item); err != nil {
			return false, fmt.Errorf("hasher Write: %w", err)
		}

		idx := int(h.Sum64() % uint64(bf.size))
		r = r && bf.bits[idx/64]&(1<<(idx%64)) != 0
	}
	return r, nil
}
