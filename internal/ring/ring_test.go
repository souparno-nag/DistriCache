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
