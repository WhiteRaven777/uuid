package uuid

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"testing"
)

func TestNewV8(t *testing.T) {
	u := NewV8(0x123456789abc, 0x0def, 0x0123456789abcdef)
	if got, want := u.String(), "12345678-9abc-8def-8123-456789abcdef"; got != want {
		t.Fatalf("NewV8() = %q, want %q", got, want)
	}
	if got := u[6] >> 4; got != 8 {
		t.Errorf("version = %d, want 8", got)
	}
	if got := u[8] >> 6; got != 0b10 {
		t.Errorf("variant = %02b, want 10", got)
	}
}

func TestNewV8TruncatesCustomFields(t *testing.T) {
	u := NewV8(^uint64(0), ^uint16(0), ^uint64(0))
	if got, want := u.String(), "ffffffff-ffff-8fff-bfff-ffffffffffff"; got != want {
		t.Fatalf("NewV8() = %q, want %q", got, want)
	}
}

func TestNewV8WithHashedRandomInput(t *testing.T) {
	var seed [32]byte
	if _, err := rand.Read(seed[:]); err != nil {
		t.Fatal(err)
	}

	customA, customB, customC := hashedV8Fields(seed[:])
	u := NewV8(customA, customB, customC)
	if got, want := binary.BigEndian.Uint64(u[0:8])>>16, customA&customAMask; got != want {
		t.Errorf("custom_a = %#x, want %#x", got, want)
	}
	if got, want := binary.BigEndian.Uint16(u[6:8])&customBMask, customB&customBMask; got != want {
		t.Errorf("custom_b = %#x, want %#x", got, want)
	}
	if got, want := binary.BigEndian.Uint64(u[8:16])&customCMask, customC&customCMask; got != want {
		t.Errorf("custom_c = %#x, want %#x", got, want)
	}
}

// hashedV8Fields demonstrates splitting a digest into UUIDv8 custom fields.
func hashedV8Fields(input []byte) (customA uint64, customB uint16, customC uint64) {
	digest := sha256.Sum256(input)
	customA = binary.BigEndian.Uint64(digest[0:8])
	customB = binary.BigEndian.Uint16(digest[8:10])
	customC = binary.BigEndian.Uint64(digest[10:18])
	return
}
