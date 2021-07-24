package uuid

import (
	"fmt"
	"testing"
	"time"
)

func TestUUIDv7(t *testing.T) {
	now := time.Now().UTC()
	id := New(UUIDv7)
	uuid := id.Convert(now)

	// short
	if id, err := Parse(uuid.Short()); err != nil {
		t.Error("parse error; short")
	} else if id.Short() != uuid.Short() {
		t.Error("mismatch; short")
	}

	// default
	if id, err := Parse(uuid.String()); err != nil {
		t.Error("parse error; default")
	} else if id.String() != uuid.String() {
		t.Error("mismatch; default")
	}

	// MS format
	if id, err := Parse(fmt.Sprintf("{%s}", uuid.String())); err != nil {
		t.Error("parse error; default")
	} else if id.String() != uuid.String() {
		t.Error("mismatch; default")
	}

	// urn
	if id, err := Parse(uuid.URN()); err != nil {
		t.Error("parse error; urn")
	} else if id.URN() != uuid.URN() {
		t.Error("mismatch; urn")
	}

	// base64
	if id, err := Parse(uuid.Base64()); err != nil {
		t.Error("parse error; base64")
	} else if id.Base64() != uuid.Base64() {
		t.Error("mismatch; base64")
	}

}

func TestVersion(t *testing.T) {
	for _, ver := range []Version{UUIDv7} {
		if ver != New(ver).Generate().Version() {
			t.Fatal("unmatched")
		}
	}
}
