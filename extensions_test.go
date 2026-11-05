package uuid

import (
	"encoding/binary"
	"testing"
	"time"
)

func TestExtensions(t *testing.T) {
	u := MustParse("01234567-89ab-7def-8123-456789abcdef")
	if got, want := u.Short(), "0123456789ab7def8123456789abcdef"; got != want {
		t.Errorf("Short() = %q, want %q", got, want)
	}
	if got, want := u.URN(), "urn:uuid:"+u.String(); got != want {
		t.Errorf("URN() = %q, want %q", got, want)
	}
	if got, want := u.Base64(), "ASNFZ4mrfe+BI0VniavN7w"; got != want {
		t.Errorf("Base64() = %q, want %q", got, want)
	}
}

func TestTime(t *testing.T) {
	const timestamp = uint64(0x1ec9414c232ab00)
	node := [6]byte{0x9f, 0x6b, 0xde, 0xce, 0xd8, 0x46}
	want := time.Date(2022, time.February, 22, 19, 22, 22, 0, time.UTC)

	for name, u := range map[string]UUID{
		"v1": formatV1(timestamp, 0x33c8, node),
		"v6": formatV6(timestamp, 0x33c8, node),
	} {
		t.Run(name, func(t *testing.T) {
			got, ok := u.Time()
			if !ok || !got.Equal(want) {
				t.Fatalf("Time() = %v, %v; want %v, true", got, ok, want)
			}
		})
	}
}

func TestV2Time(t *testing.T) {
	const timestamp = uint64(0x1ec9414c232ab00)
	u := formatV2(formatV1(timestamp, 0x33c8, [6]byte{}), Person, 1000)
	got, ok := u.Time()
	want := gregorianTime(timestamp &^ 0xffffffff)
	if !ok || !got.Equal(want) {
		t.Fatalf("Time() = %v, %v; want %v, true", got, ok, want)
	}
}

func TestV7Time(t *testing.T) {
	const milliseconds = int64(1_645_557_742_000)
	var u UUID
	var encoded [8]byte
	binary.BigEndian.PutUint64(encoded[:], uint64(milliseconds))
	copy(u[0:6], encoded[2:])
	u[6] = 0x70
	u[8] = 0x80

	got, ok := u.Time()
	want := time.UnixMilli(milliseconds).UTC()
	if !ok || !got.Equal(want) {
		t.Fatalf("Time() = %v, %v; want %v, true", got, ok, want)
	}
}

func TestTimeUnavailable(t *testing.T) {
	for _, u := range []UUID{
		NewV3(DNS, []byte("example.com")),
		NewV4(),
		NewV5(DNS, []byte("example.com")),
		NewV8(1, 2, 3),
	} {
		if got, ok := u.Time(); ok || !got.IsZero() {
			t.Errorf("Time() = %v, %v; want zero time, false", got, ok)
		}
	}
}
