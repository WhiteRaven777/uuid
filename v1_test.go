package uuid

import (
	"encoding/binary"
	"net"
	"sync"
	"testing"
	"time"
)

func TestV1RFC9562Vector(t *testing.T) {
	node := [6]byte{0x9f, 0x6b, 0xde, 0xce, 0xd8, 0x46}
	u := formatV1(0x1ec9414c232ab00, 0x33c8, node)
	if got, want := u.String(), "c232ab00-9414-11ec-b3c8-9f6bdeced846"; got != want {
		t.Fatalf("formatV1() = %q, want %q", got, want)
	}
}

func TestGregorianTimestamp(t *testing.T) {
	timestamp := time.Date(2022, time.February, 22, 19, 22, 22, 0, time.UTC)
	if got, want := gregorianTimestamp(timestamp), uint64(0x1ec9414c232ab00); got != want {
		t.Fatalf("gregorianTimestamp() = %#x, want %#x", got, want)
	}
}

func TestNewV1(t *testing.T) {
	u := NewV1()
	if got := u[6] >> 4; got != 1 {
		t.Errorf("version = %d, want 1", got)
	}
	if got := u[8] >> 6; got != 0b10 {
		t.Errorf("variant = %02b, want 10", got)
	}
	if got := u[10] & 0x01; got != 1 {
		t.Errorf("node multicast bit = %d, want 1", got)
	}
}

func TestV1ClockRollback(t *testing.T) {
	g := v1Generator{
		initialised:   true,
		node:          [6]byte{1, 2, 3, 4, 5, 6},
		clockSequence: 7,
	}
	now := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	first := g.new(now, nil)
	second := g.new(now.Add(-time.Second), nil)

	firstSequence := binary.BigEndian.Uint16(first[8:10]) & 0x3fff
	secondSequence := binary.BigEndian.Uint16(second[8:10]) & 0x3fff
	if got, want := secondSequence, (firstSequence+1)&0x3fff; got != want {
		t.Errorf("clock sequence after rollback = %d, want %d", got, want)
	}
	if first == second {
		t.Fatal("clock rollback generated duplicate UUIDs")
	}
}

func TestV1MACNode(t *testing.T) {
	want := [6]byte{0x02, 0x11, 0x22, 0x33, 0x44, 0x55}
	g := v1Generator{}
	u := g.new(time.Now(), func() ([6]byte, bool) {
		return want, true
	})
	if got := [6]byte(u[10:]); got != want {
		t.Fatalf("node = %x, want %x", got, want)
	}
}

func TestV1MACNodeFallback(t *testing.T) {
	g := v1Generator{}
	u := g.new(time.Now(), func() ([6]byte, bool) {
		return [6]byte{}, false
	})
	if got := u[10] & 0x01; got != 1 {
		t.Fatalf("fallback node multicast bit = %d, want 1", got)
	}
}

func TestMACNode(t *testing.T) {
	want := [6]byte{0x02, 0x11, 0x22, 0x33, 0x44, 0x55}
	got, ok := macNode(net.HardwareAddr(want[:]))
	if !ok || got != want {
		t.Fatalf("macNode() = %x, %v; want %x, true", got, ok, want)
	}

	for _, address := range []net.HardwareAddr{
		nil,
		{0, 0, 0, 0, 0, 0},
		{0x01, 0x11, 0x22, 0x33, 0x44, 0x55},
		{0x02, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77},
	} {
		if got, ok := macNode(address); ok {
			t.Errorf("macNode(%x) = %x, true; want false", address, got)
		}
	}
}

func TestNewV1Concurrent(t *testing.T) {
	const (
		workers = 8
		perWork = 2_000
	)

	values := make(chan UUID, workers*perWork)
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range perWork {
				values <- NewV1()
			}
		}()
	}
	wg.Wait()
	close(values)

	seen := make(map[UUID]struct{}, workers*perWork)
	for u := range values {
		if _, exists := seen[u]; exists {
			t.Fatalf("NewV1 generated duplicate UUID %v", u)
		}
		seen[u] = struct{}{}
	}
}
