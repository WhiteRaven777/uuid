package uuid

import "testing"

func TestV6RFC9562Vector(t *testing.T) {
	node := [6]byte{0x9f, 0x6b, 0xde, 0xce, 0xd8, 0x46}
	u := formatV6(0x1ec9414c232ab00, 0x33c8, node)
	if got, want := u.String(), "1ec9414c-232a-6b00-b3c8-9f6bdeced846"; got != want {
		t.Fatalf("formatV6() = %q, want %q", got, want)
	}
}

func TestNewV6(t *testing.T) {
	u := NewV6()
	if got := u[6] >> 4; got != 6 {
		t.Errorf("version = %d, want 6", got)
	}
	if got := u[8] >> 6; got != 0b10 {
		t.Errorf("variant = %02b, want 10", got)
	}
	if got := u[10] & 0x01; got != 1 {
		t.Errorf("node multicast bit = %d, want 1", got)
	}
}

func TestNewV6WithMAC(t *testing.T) {
	u := NewV6(true)
	if got := u[6] >> 4; got != 6 {
		t.Errorf("version = %d, want 6", got)
	}
	if got := u[8] >> 6; got != 0b10 {
		t.Errorf("variant = %02b, want 10", got)
	}
}

func TestNewV6SortOrder(t *testing.T) {
	previous := NewV6()
	for range 10_000 {
		current := NewV6()
		if got := current.Compare(previous); got != 1 {
			t.Fatalf("NewV6 returned UUIDs out of order:\nprevious: %v\n current: %v", previous, current)
		}
		previous = current
	}
}
