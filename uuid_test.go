package uuid_test

import (
	"encoding"
	"testing"

	"github.com/WhiteRaven777/uuid"
)

func TestNew(t *testing.T) {
	for _, test := range []struct {
		name    string
		newf    func() uuid.UUID
		version byte
		variant byte
	}{
		{name: "New", newf: uuid.New, version: 4, variant: 0b10},
		{name: "NewV4", newf: uuid.NewV4, version: 4, variant: 0b10},
		{name: "NewV7", newf: uuid.NewV7, version: 7, variant: 0b10},
	} {
		u := test.newf()
		if got := u[6] >> 4; got != test.version {
			t.Errorf("%s: got version %d, want %d", test.name, got, test.version)
		}
		if got := u[8] >> 6; got != test.variant {
			t.Errorf("%s: got variant %02b, want %02b", test.name, got, test.variant)
		}
	}
}

func TestEncode(t *testing.T) {
	u := uuid.UUID{0xf8, 0x1d, 0x4f, 0xae, 0x7d, 0xec, 0x11, 0xd0, 0xa7, 0x65, 0x00, 0xa0, 0xc9, 0x1e, 0x6b, 0xf6}
	want := "f81d4fae-7dec-11d0-a765-00a0c91e6bf6"
	if got := u.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
	if got, err := u.MarshalText(); string(got) != want || err != nil {
		t.Errorf("MarshalText() = %q, %v; want %q, nil", got, err, want)
	}
	prefix := []byte("uuid=")
	if got, err := u.AppendText(prefix); string(got) != string(prefix)+want || err != nil {
		t.Errorf("AppendText(%q) = %q, %v; want %q, nil", prefix, got, err, string(prefix)+want)
	}
}

func TestUnmarshalText(t *testing.T) {
	var got uuid.UUID
	if err := got.UnmarshalText([]byte("f81d4fae-7dec-11d0-a765-00a0c91e6bf6")); err != nil {
		t.Fatalf("UnmarshalText() returned %v", err)
	}
	want := uuid.UUID{0xf8, 0x1d, 0x4f, 0xae, 0x7d, 0xec, 0x11, 0xd0, 0xa7, 0x65, 0x00, 0xa0, 0xc9, 0x1e, 0x6b, 0xf6}
	if got != want {
		t.Errorf("UnmarshalText() = %v, want %v", got, want)
	}
}

func TestParseSuccess(t *testing.T) {
	want := uuid.UUID{0xf8, 0x1d, 0x4f, 0xae, 0x7d, 0xec, 0x11, 0xd0, 0xa7, 0x65, 0x00, 0xa0, 0xc9, 0x1e, 0x6b, 0xf6}
	for _, input := range []string{
		"f81d4fae-7dec-11d0-a765-00a0c91e6bf6",
		"F81D4FAE-7DEC-11D0-A765-00A0C91E6BF6",
		"f81d4fae7dec11d0a76500a0c91e6bf6",
		"{f81d4fae-7dec-11d0-a765-00a0c91e6bf6}",
		"urn:uuid:f81d4fae-7dec-11d0-a765-00a0c91e6bf6",
	} {
		got, err := uuid.Parse(input)
		if err != nil {
			t.Errorf("Parse(%q) returned %v", input, err)
		} else if got != want {
			t.Errorf("Parse(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestParseErrors(t *testing.T) {
	for _, input := range []string{
		"",
		"not-a-uuid",
		"00000000-0000-0000-0000-00000000000",
		"00000000-0000-0000-0000-00000000000x",
	} {
		if got, err := uuid.Parse(input); err == nil {
			t.Errorf("Parse(%q) = %v, nil; want error", input, got)
		}
	}
}

func TestCompare(t *testing.T) {
	uuids := []uuid.UUID{
		uuid.Nil(),
		uuid.MustParse("f81d4fae-7dec-11d0-a765-00a0c91e6bf6"),
		uuid.Max(),
	}
	for i, u := range uuids {
		if got := u.Compare(u); got != 0 {
			t.Errorf("%v.Compare(itself) = %d, want 0", u, got)
		}
		if i == 0 {
			continue
		}
		previous := uuids[i-1]
		if got := previous.Compare(u); got != -1 {
			t.Errorf("%v.Compare(%v) = %d, want -1", previous, u, got)
		}
		if got := u.Compare(previous); got != 1 {
			t.Errorf("%v.Compare(%v) = %d, want 1", u, previous, got)
		}
	}
}

func TestEncodingInterfaces(t *testing.T) {
	var _ encoding.TextMarshaler = uuid.UUID{}
	var _ encoding.TextUnmarshaler = (*uuid.UUID)(nil)
}
