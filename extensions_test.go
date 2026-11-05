package uuid_test

import (
	"testing"

	"github.com/WhiteRaven777/uuid"
)

func TestExtensions(t *testing.T) {
	u := uuid.MustParse("01234567-89ab-7def-8123-456789abcdef")
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
