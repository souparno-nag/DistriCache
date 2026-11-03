package ring

import (
	"hash/fnv"
	"sort"
	"strconv"
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

func (r *Ring) AddNode(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.nodes[id] {
		return
	}
	r.nodes[id] = true

	for i := 0; i < r.vnodes; i++ {
		h := hash(id + "#" + strconv.Itoa(i))
		r.hashes = append(r.hashes, h)
		r.owner[h] = id
	}

	sort.Slice(r.hashes, func(i, j int) bool { return r.hashes[i] < r.hashes[j] })
}
