package bf

type BloomFilter interface {
	Add(item []byte) error
	Contains(item []byte) (bool, error)
}
