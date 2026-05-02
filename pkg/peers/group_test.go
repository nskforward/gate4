package peers

import (
	"context"
	"sync"
	"testing"
)

func TestNewGroup(t *testing.T) {
	ps := NewPubSub[string]()
	group := NewGroup("test-key", ps)

	if group == nil {
		t.Fatal("expected non-nil group")
	}
	if group.key != "test-key" {
		t.Errorf("expected key 'test-key', got %q", group.key)
	}
	if group.pubsub != ps {
		t.Error("expected pubsub to match")
	}
	if group.last == nil {
		t.Error("expected last to be non-nil")
	}
}

func TestGroupNewPeer(t *testing.T) {
	ps := NewPubSub[string]()
	group := NewGroup("test", ps)

	peer := group.NewPeer()

	if peer == nil {
		t.Fatal("expected non-nil peer")
	}
	if len(group.peers) != 1 {
		t.Errorf("expected 1 peer in group, got %d", len(group.peers))
	}

	// Create another peer
	peer2 := group.NewPeer()
	if len(group.peers) != 2 {
		t.Errorf("expected 2 peers in group, got %d", len(group.peers))
	}
	if peer == peer2 {
		t.Error("expected different peer instances")
	}
}

func TestGroupNewPeerGetsLastValue(t *testing.T) {
	ps := NewPubSub[string]()
	group := NewGroup("test", ps)

	// Send a value first
	group.Send("initial")

	// Create a new peer - it should receive the last value
	peer := group.NewPeer()

	ctx := context.Background()
	data, ok := peer.Read(ctx)

	if !ok {
		t.Fatal("expected to receive last value")
	}
	if data != "initial" {
		t.Errorf("expected 'initial', got %q", data)
	}
}

func TestGroupSend(t *testing.T) {
	ps := NewPubSub[string]()
	group := NewGroup("test", ps)

	peer1 := group.NewPeer()
	peer2 := group.NewPeer()

	// Send to all peers
	result := group.Send("message")

	if !result {
		t.Error("expected Send to return true")
	}

	// Both peers should receive the message
	ctx := context.Background()

	data1, ok1 := peer1.Read(ctx)
	if !ok1 {
		t.Error("expected peer1 to receive message")
	}
	if data1 != "message" {
		t.Errorf("expected 'message', got %q", data1)
	}

	data2, ok2 := peer2.Read(ctx)
	if !ok2 {
		t.Error("expected peer2 to receive message")
	}
	if data2 != "message" {
		t.Errorf("expected 'message', got %q", data2)
	}
}

func TestGroupSendNoPeers(t *testing.T) {
	ps := NewPubSub[string]()
	group := NewGroup("test", ps)

	result := group.Send("message")

	if result {
		t.Error("expected Send to return false when no peers")
	}
}

func TestGroupLastValue(t *testing.T) {
	ps := NewPubSub[string]()
	group := NewGroup("test", ps)

	group.Send("first")
	group.Send("second")

	// Create a new peer, it should get the last value ("second")
	peer := group.NewPeer()

	ctx := context.Background()
	data, ok := peer.Read(ctx)

	if !ok {
		t.Fatal("expected to receive last value")
	}
	if data != "second" {
		t.Errorf("expected 'second', got %q", data)
	}
}

func TestGroupClose(t *testing.T) {
	ps := NewPubSub[string]()
	group := NewGroup("test-key", ps)

	_ = group.NewPeer()
	_ = group.NewPeer()

	if len(group.peers) != 2 {
		t.Errorf("expected 2 peers, got %d", len(group.peers))
	}

	group.Close()

	// Peers should be closed and removed
	if len(group.peers) != 0 {
		t.Errorf("expected 0 peers after close, got %d", len(group.peers))
	}

	// Group should be removed from pubsub
	group2, loaded := ps.LoadOrCreate("test-key")
	if loaded {
		t.Error("expected group to be removed from pubsub")
	}
	// A new group should be created
	if group2 == group {
		t.Error("expected new group to be different")
	}
}

func TestGroupSendConcurrency(t *testing.T) {
	ps := NewPubSub[string]()
	group := NewGroup("test", ps)

	// Create multiple peers
	for i := 0; i < 10; i++ {
		group.NewPeer()
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			group.Send(string(rune('a' + n%26)))
		}(i)
	}

	wg.Wait()

	// All peers should have received some messages
	for i, peer := range group.peers {
		if len(peer.channel) == 0 {
			t.Errorf("peer %d should have received messages", i)
		}
	}
}

func TestGroupPeerRemoval(t *testing.T) {
	ps := NewPubSub[string]()
	group := NewGroup("test", ps)

	peer1 := group.NewPeer()
	peer2 := group.NewPeer()

	if len(group.peers) != 2 {
		t.Errorf("expected 2 peers, got %d", len(group.peers))
	}

	// Close one peer
	peer1.Close()

	if len(group.peers) != 1 {
		t.Errorf("expected 1 peer after close, got %d", len(group.peers))
	}

	// The remaining peer should still work
	group.Send("after removal")

	ctx := context.Background()
	data, ok := peer2.Read(ctx)
	if !ok {
		t.Error("expected peer2 to receive message")
	}
	if data != "after removal" {
		t.Errorf("expected 'after removal', got %q", data)
	}
}
