package ring

import (
	"fmt"
	"sync"
	"testing"
)

func TestSameKeySameNode(t *testing.T) {
	r := New(150)
	r.AddNode("node-0")
	r.AddNode("node-1")
	r.AddNode("node-2")

	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key-%d", i)
		first, ok := r.GetNode(key)
		if !ok {
			t.Fatalf("GetNode(%q) found no node on a non-empty ring", key)
		}
		for j := 0; j < 5; j++ {
			if got, _ := r.GetNode(key); got != first {
				t.Fatalf("key %q moved from %s to %s", key, first, got)
			}
		}
	}
}

func TestEmptyRing(t *testing.T) {
	r := New(150)
	node, ok := r.GetNode("user:42")
	if ok {
		t.Fatalf("empty ring returned node %q, want ok=false", node)
	}
}

func TestAddNodeTwice(t *testing.T) {
	r := New(150)
	r.AddNode("node-1")
	before := len(r.hashes)

	r.AddNode("node-1")

	if after := len(r.hashes); after != before {
		t.Fatalf("adding node twice changed positions from %d to %d", before, after)
	}
}

func TestRemoveNode(t *testing.T) {
	r := New(150)
	r.AddNode("node-0")
	r.AddNode("node-1")
	r.AddNode("node-2")

	r.RemoveNode("node-1")

	if len(r.hashes) != 300 {
		t.Fatalf("len(r.hashes) = %d, want 300", len(r.hashes))
	}

	if len(r.nodes) != 2 {
		t.Fatalf("len(r.nodes) = %d, want 2", len(r.nodes))
	}

	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%d", i)
		node, ok := r.GetNode(key)

		if !ok {
			t.Fatalf("GetNode(%q) found no node", key)
		}

		if node == "node-1" {
			t.Fatalf("key %q still routes to removed node-1", key)
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	r := New(150)
	r.AddNode("node-0")
	r.AddNode("node-2")

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				key := fmt.Sprintf("key-%d-%d", i, j)
				if _, ok := r.GetNode(key); !ok {
					t.Errorf("GetNode(%q) returned ok=false", key)
				}
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 100; j++ {
			r.AddNode("node-1")
			r.RemoveNode("node-1")
		}
	}()

	wg.Wait()
}

func TestGetNodesFirstMatchesGetNode(t *testing.T) {
	r := New(150)
	r.AddNode("node-0")
	r.AddNode("node-1")
	r.AddNode("node-2")

	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key-%d", i)
		first, ok := r.GetNode(key)
		if !ok {
			t.Fatalf("GetNode(%q) found no node on a non-empty ring", key)
		}
		check := r.GetNodes(key, 1)
		if len(check) != 1 {
			t.Fatalf("GetNodes(%q, 1) returned %d nodes, want 1", key, len(check))
		}
		if check[0] != first {
			t.Fatalf("GetNodes(%q, 1)[0] = %s, want %s (from GetNode)", key, check[0], first)
		}
	}
}
