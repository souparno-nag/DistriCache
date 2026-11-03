package ring

import (
	"fmt"
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
