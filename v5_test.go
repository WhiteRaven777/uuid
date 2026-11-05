package uuid_test

import (
	"testing"

	"github.com/WhiteRaven777/uuid"
)

func TestNewV5RFC9562Vector(t *testing.T) {
	u := uuid.NewV5(uuid.DNS, []byte("www.example.com"))
	if got, want := u.String(), "2ed6657d-e927-568b-95e1-2665a8aea6a2"; got != want {
		t.Fatalf("NewV5() = %q, want %q", got, want)
	}
}

func TestNewV5(t *testing.T) {
	name := []byte("example.com")
	u := uuid.NewV5(uuid.DNS, name)
	if got := u[6] >> 4; got != 5 {
		t.Errorf("version = %d, want 5", got)
	}
	if got := u[8] >> 6; got != 0b10 {
		t.Errorf("variant = %02b, want 10", got)
	}
	if got := uuid.NewV5(uuid.DNS, name); got != u {
		t.Errorf("same namespace and name produced %v, want %v", got, u)
	}
	if got := uuid.NewV5(uuid.URL, name); got == u {
		t.Errorf("different namespaces produced the same UUID %v", got)
	}
	if got := uuid.NewV5(uuid.DNS, []byte("example.net")); got == u {
		t.Errorf("different names produced the same UUID %v", got)
	}
}

func TestNewV5InvalidNamespace(t *testing.T) {
	if got := uuid.NewV5(uuid.Namespace(0), []byte("example.com")); got != uuid.Nil() {
		t.Fatalf("NewV5() = %v for an unsupported namespace, want Nil UUID", got)
	}
}
