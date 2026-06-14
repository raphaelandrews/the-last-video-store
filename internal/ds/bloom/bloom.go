package bloom

import "hash/fnv"

type BloomFilter struct {
	bitset    []uint64
	size      uint64
	hashCount int
}

func New(size uint64, hashCount int) *BloomFilter {
	if size < 1 {
		size = 1024
	}
	if hashCount < 1 {
		hashCount = 3
	}
	arrLen := (size + 63) / 64
	return &BloomFilter{
		bitset:    make([]uint64, arrLen),
		size:      size,
		hashCount: hashCount,
	}
}

func (bf *BloomFilter) Add(data []byte) {
	h1, h2 := bf.hashes(data)
	for i := range bf.hashCount {
		idx := (h1 + uint64(i)*h2) % bf.size
		wordIdx := idx / 64
		bitIdx := idx % 64
		bf.bitset[wordIdx] |= 1 << bitIdx
	}
}

func (bf *BloomFilter) Contains(data []byte) bool {
	h1, h2 := bf.hashes(data)
	for i := range bf.hashCount {
		idx := (h1 + uint64(i)*h2) % bf.size
		wordIdx := idx / 64
		bitIdx := idx % 64
		if bf.bitset[wordIdx]&(1<<bitIdx) == 0 {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) hashes(data []byte) (uint64, uint64) {
	h1 := fnv.New64a()
	h1.Write(data)
	hash1 := h1.Sum64()

	h2 := fnv.New64()
	h2.Write(data)
	hash2 := h2.Sum64()

	if hash2 == 0 {
		hash2 = 1
	}
	return hash1, hash2
}
