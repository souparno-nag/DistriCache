package ring

import (
	"hash/fnv"
	"sync"
)

type Ring struct {
	mu     sync.RWMutex
	vnodes int
	hashes []uint32
	owner  map[uint32]string
	nodes  map[string]bool
}

func New(vnodes int) *Ring {
	return &Ring{
		vnodes: vnodes,
		owner:  make(map[uint32]string),
		nodes:  make(map[string]bool),
	}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
