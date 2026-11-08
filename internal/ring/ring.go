package ring

import (
	"hash/fnv"
	"slices"
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
		name := id + "#" + strconv.Itoa(i)
		h := hash(name)

		for attempt := 1; ; attempt++ {
			if _, taken := r.owner[h]; !taken {
				break
			}
			h = hash(name + "#" + strconv.Itoa(attempt))
		}

		r.hashes = append(r.hashes, h)
		r.owner[h] = id
	}

	sort.Slice(r.hashes, func(i, j int) bool { return r.hashes[i] < r.hashes[j] })
}

func (r *Ring) GetNode(key string) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.hashes) == 0 {
		return "", false
	}

	h := hash(key)

	index := sort.Search(len(r.hashes), func(i int) bool { return r.hashes[i] >= h })

	if index == len(r.hashes) {
		index = 0
	}

	return r.owner[r.hashes[index]], true
}

func (r *Ring) GetNodes(key string, n int) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.hashes) == 0 || n <= 0 {
		return nil
	}

	h := hash(key)

	index := sort.Search(len(r.hashes), func(i int) bool { return r.hashes[i] >= h })
	if index == len(r.hashes) {
		index = 0
	}

	result := make([]string, 0, n)
	seen := make(map[string]bool)

	i := index
	for step := 0; step < len(r.hashes) && len(result) < n; step++ {
		node := r.owner[r.hashes[i]]
		if !seen[node] {
			seen[node] = true
			result = append(result, node)
		}
		i = (i + 1) % len(r.hashes)
	}

	return result
}

func (r *Ring) RemoveNode(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.nodes[id] {
		return
	}
	delete(r.nodes, id)

	r.hashes = slices.DeleteFunc(r.hashes, func(h uint32) bool {
		if r.owner[h] == id {
			delete(r.owner, h)
			return true
		}
		return false
	})
}
