package uuid

import (
	"encoding/binary"
	"testing"
	"time"
)

func TestFormatV2(t *testing.T) {
	node := [6]byte{0x9f, 0x6b, 0xde, 0xce, 0xd8, 0x46}
	u := formatV2(formatV1(0x1ec9414c232ab00, 0x33c8, node), Group, 1000)
	if got, want := u.String(), "000003e8-9414-21ec-b301-9f6bdeced846"; got != want {
		t.Fatalf("formatV2() = %q, want %q", got, want)
	}
}

func TestNewV2(t *testing.T) {
	const localID = uint32(1000)
	u := NewV2(Person, localID)
	if got := binary.BigEndian.Uint32(u[0:4]); got != localID {
		t.Errorf("local ID = %d, want %d", got, localID)
	}
	if got := u[6] >> 4; got != 2 {
		t.Errorf("version = %d, want 2", got)
	}
	if got := u[8] >> 6; got != 0b10 {
		t.Errorf("variant = %02b, want 10", got)
	}
	if got := LocalDomain(u[9]); got != Person {
		t.Errorf("local domain = %d, want %d", got, Person)
	}
	if got := u[10] & 0x01; got != 1 {
		t.Errorf("node multicast bit = %d, want 1", got)
	}
}

func TestNewV2Domains(t *testing.T) {
	for _, domain := range []LocalDomain{Person, Group, Organisation} {
		u := NewV2(domain, 1000)
		if got := LocalDomain(u[9]); got != domain {
			t.Errorf("NewV2(%d) domain = %d, want %d", domain, got, domain)
		}
	}
}

func TestNewV2WithMAC(t *testing.T) {
	u := NewV2(Group, 1000, true)
	if got := u[6] >> 4; got != 2 {
		t.Errorf("version = %d, want 2", got)
	}
	if got := LocalDomain(u[9]); got != Group {
		t.Errorf("local domain = %d, want %d", got, Group)
	}
}

func TestNewV2InvalidDomain(t *testing.T) {
	if got := NewV2(LocalDomain(3), 1000); got != Nil() {
		t.Fatalf("NewV2() = %v for an unsupported local domain, want Nil UUID", got)
	}
}

func TestDCEGeneratorClockSequence(t *testing.T) {
	g := dceGenerator{}
	now := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	seen := make(map[uint16]struct{}, 64)
	for range 64 {
		_, clockSequence, _ := g.next(now, nil)
		if _, exists := seen[clockSequence]; exists {
			t.Fatalf("clock sequence %d repeated before all 64 values were used", clockSequence)
		}
		seen[clockSequence] = struct{}{}
	}
}
