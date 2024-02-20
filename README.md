# go-bf

go-bf is a Go package that provides an implementation of a Bloom filter. A Bloom filter is a space-efficient probabilistic data structure that is used to test whether an element is a member of a set. It is capable of telling you if an element might be in the set or it is definitely not in the set.

## Features
* Customizable Hash Functions: The package allows you to specify custom hash functions or use the default FNV-1a hash function.
* Configurable Hash Count: You can configure the number of hash functions to use for the Bloom filter.
* Concurrency-Safe(???): The package is designed to be safe for concurrent use, but you should ensure that the hash functions you provide are also safe for concurrent use.

## Installation
To install the package, run:
```shell
go get github.com/robotomize/go-bf/bf/bfbits
```

## Usage
To use bfbits, you need to import the package and create a new Bloom filter with the desired size and options.

```go
// Create a new Bloom filter with the default settings.
bf := bfbits.NewBloomFilter(10000)

// Add an item to the Bloom filter.
err := bf.Add([]byte("apple"))
if err != nil {
    log.Fatal(err)
}

// Check if an item is in the Bloom filter.
contains, err := bf.Contains([]byte("apple"))
if err != nil {
    log.Fatal(err)
}
fmt.Printf("item exist: %t", contains)
```
## Options
You can customize the behavior of the Bloom filter by providing options when creating it.
```go
// Create a Bloom filter with a custom hash function and hash count.
customHashFunc := func() hash.Hash64 {
    return fnv.New64a()
}
bf := bfbits.NewBloomFilter(10000, bfbits.WithHasher(customHashFunc), bfbits.WithHashNum(2))
```
## License
This project is licensed under the MIT License. See the LICENSE file for details.