package ring

import "sync"

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
