package peers

import (
	"context"
	"testing"
	"time"
)

func TestNewPeer(t *testing.T) {
	ps := NewPubSub[string]()
	group := NewGroup("test", ps)

	peer := NewPeer(8, group)

	if peer == nil {
		t.Fatal("expected non-nil peer")
	}
	if peer.group != group {
		t.Error("expected group to match")
	}
	if peer.channel == nil {
		t.Error("expected channel to be non-nil")
	}
	if cap(peer.channel) != 8 {
		t.Errorf("expected channel capacity 8, got %d", cap(peer.channel))
	}
}

func TestPeerRead(t *testing.T) {
	ps := NewPubSub[string]()
	group := NewGroup("test", ps)
	peer := group.NewPeer()

	// Send data via send method (bypassing channel directly for setup)
	peer.send("hello")

	// Now read it
	ctx := context.Background()
	data, ok := peer.Read(ctx)

	if !ok {
		t.Error("expected ok to be true")
	}
	if data != "hello" {
		t.Errorf("expected 'hello', got %q", data)
	}
}

func TestPeerReadCancel(t *testing.T) {
	ps := NewPubSub[string]()
	group := NewGroup("test", ps)
	peer := group.NewPeer()

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	data, ok := peer.Read(ctx)

	if ok {
		t.Error("expected ok to be false for cancelled context")
	}
	if data != "" {
		t.Errorf("expected empty string, got %q", data)
	}
}

func TestPeerClose(t *testing.T) {
	ps := NewPubSub[string]()
	group := NewGroup("test", ps)
	peer := group.NewPeer()

	// Verify peer is in group
	if len(group.peers) != 1 {
		t.Errorf("expected 1 peer in group, got %d", len(group.peers))
	}

	peer.Close()

	// Verify peer is removed from group
	if len(group.peers) != 0 {
		t.Errorf("expected 0 peers after close, got %d", len(group.peers))
	}

	// Trying to read from closed channel should return false
	ctx := context.Background()
	_, ok := peer.Read(ctx)
	if ok {
		t.Error("expected ok to be false for closed channel")
	}
}

func TestPeerReadEmpty(t *testing.T) {
	ps := NewPubSub[string]()
	group := NewGroup("test", ps)
	peer := group.NewPeer()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	data, ok := peer.Read(ctx)

	if ok {
		t.Error("expected ok to be false for empty channel with timeout")
	}
	if data != "" {
		t.Errorf("expected empty string, got %q", data)
	}
}
