package uuid

import (
	"strings"
	"testing"
	"time"
)

func TestV7(t *testing.T) {
	id := New(UUIDv7)
	var buf UUID
	old := make(map[UUID]struct{})
	for _, precision := range []Precision{PrecisionDefault, PrecisionLow, PrecisionMid, PrecisionHigh} {
		id.(*V7).SetPrecision(precision)
		for i := 0; i < (1 << 20); i++ {
			buf = id.Generate()
			if _, ok := old[buf]; ok {
				t.Fatal("collision!!")
			}
			old[buf] = struct{}{}
		}
	}

	now := time.Now().UTC()
	for _, precision := range []Precision{PrecisionDefault, PrecisionLow, PrecisionMid, PrecisionHigh} {
		id.(*V7).SetPrecision(precision)
		for i := 0; i < (1 << 20); i++ {
			if now.Unix() != id.Decode(id.Convert(now)).Unix() {
				t.Fatal("miss decode")
			}
		}
	}
}

func TestMonotonicClock(t *testing.T) {
	var buf uint64
	old := make(map[uint64]struct{})
	for i := uint64(0); i < (1 << 20); i++ {
		buf = getMonotonicClock(time.Now())
		if _, ok := old[buf]; ok {
			t.Fatal("collision!", buf)
		}
		old[buf] = struct{}{}
	}
}

func TestS(t *testing.T) {
	f := `([0-9a-z]{8})-([0-9a-z]{4})-([0-9a-z]{4})-([0-9a-z]{4})-([0-9a-z]{12})`
	t.Log(strings.ReplaceAll(f, ")-(", ")("))
}
