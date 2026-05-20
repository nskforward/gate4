package peers

import (
	"sync"
	"testing"
)

func TestNewPubSub(t *testing.T) {
	ps := NewPubSub[string]()

	if ps == nil {
		t.Fatal("expected non-nil pubsub")
	}
	if ps.groups == nil {
		t.Error("expected groups map to be non-nil")
	}
	if len(ps.groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(ps.groups))
	}
}

func TestPubSubLoadOrCreateNew(t *testing.T) {
	ps := NewPubSub[string]()

	group, loaded := ps.LoadOrCreate("key1")

	if loaded {
		t.Error("expected loaded to be false for new group")
	}
	if group == nil {
		t.Fatal("expected non-nil group")
	}
	if group.key != "key1" {
		t.Errorf("expected key 'key1', got %q", group.key)
	}
	if len(ps.groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(ps.groups))
	}
}

func TestPubSubLoadOrCreateExisting(t *testing.T) {
	ps := NewPubSub[string]()

	group1, loaded1 := ps.LoadOrCreate("key1")
	if loaded1 {
		t.Error("expected loaded to be false for first creation")
	}

	group2, loaded2 := ps.LoadOrCreate("key1")

	if !loaded2 {
		t.Error("expected loaded to be true for existing group")
	}
	if group2 != group1 {
		t.Error("expected same group instance")
	}
	if len(ps.groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(ps.groups))
	}
}

func TestPubSubLoadOrCreateMultipleKeys(t *testing.T) {
	ps := NewPubSub[string]()

	group1, _ := ps.LoadOrCreate("key1")
	group2, _ := ps.LoadOrCreate("key2")
	group3, _ := ps.LoadOrCreate("key3")

	if len(ps.groups) != 3 {
		t.Errorf("expected 3 groups, got %d", len(ps.groups))
	}

	if group1 == group2 || group2 == group3 || group1 == group3 {
		t.Error("expected different group instances")
	}
}

func TestPubSubLoadOrCreateConcurrency(t *testing.T) {
	ps := NewPubSub[string]()

	var wg sync.WaitGroup
	results := make(chan *Group[string], 100)

	// Create multiple groups concurrently with same key
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			group, _ := ps.LoadOrCreate("same-key")
			results <- group
		}()
	}

	wg.Wait()
	close(results)

	// All results should be the same instance
	var first *Group[string]
	for group := range results {
		if first == nil {
			first = group
		} else if group != first {
			t.Error("expected all groups to be the same instance")
		}
	}

	if len(ps.groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(ps.groups))
	}
}

func TestPubSubRemove(t *testing.T) {
	ps := NewPubSub[string]()

	group, _ := ps.LoadOrCreate("key1")

	if len(ps.groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(ps.groups))
	}

	// Remove should be called internally by Group.Close
	group.Close()

	if len(ps.groups) != 0 {
		t.Errorf("expected 0 groups after removal, got %d", len(ps.groups))
	}

	// Create a new group with same key
	group2, loaded := ps.LoadOrCreate("key1")

	if loaded {
		t.Error("expected loaded to be false for new group after removal")
	}
	if group2 == group {
		t.Error("expected new group to be different instance")
	}
}

func TestPubSubWithDifferentTypes(t *testing.T) {
	psInt := NewPubSub[int]()
	psStr := NewPubSub[string]()
	psBool := NewPubSub[bool]()

	groupInt, _ := psInt.LoadOrCreate("key")
	groupStr, _ := psStr.LoadOrCreate("key")
	groupBool, _ := psBool.LoadOrCreate("key")

	// They should work independently
	groupInt.Send(42)
	groupStr.Send("test")
	groupBool.Send(true)

	// Verify data types are correct
	peerInt := groupInt.NewPeer()
	peerStr := groupStr.NewPeer()
	peerBool := groupBool.NewPeer()

	// Check channel types
	intChan := peerInt.channel
	strChan := peerStr.channel
	boolChan := peerBool.channel

	if cap(intChan) == 0 {
		t.Error("expected int channel to be valid")
	}
	if cap(strChan) == 0 {
		t.Error("expected string channel to be valid")
	}
	if cap(boolChan) == 0 {
		t.Error("expected bool channel to be valid")
	}
}
