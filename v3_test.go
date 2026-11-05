package uuid_test

import (
	"testing"

	"github.com/WhiteRaven777/uuid"
)

func TestNamespaces(t *testing.T) {
	namespaces := []uuid.Namespace{
		uuid.DNS,
		uuid.URL,
		uuid.OID,
		uuid.X500,
	}
	seen := make(map[uuid.Namespace]struct{}, len(namespaces))
	generated := make(map[uuid.UUID]struct{}, len(namespaces))
	for _, namespace := range namespaces {
		if _, exists := seen[namespace]; exists {
			t.Errorf("duplicate namespace mode %d", namespace)
		}
		seen[namespace] = struct{}{}

		u := uuid.NewV3(namespace, []byte("www.example.com"))
		if u == uuid.Nil() {
			t.Errorf("NewV3(%d) returned Nil UUID", namespace)
		}
		if _, exists := generated[u]; exists {
			t.Errorf("namespace %d generated duplicate UUID %v", namespace, u)
		}
		generated[u] = struct{}{}
	}
}

func TestNewV3RFC9562Vector(t *testing.T) {
	u := uuid.NewV3(uuid.DNS, []byte("www.example.com"))
	if got, want := u.String(), "5df41881-3aed-3515-88a7-2f4a814cf09e"; got != want {
		t.Fatalf("NewV3() = %q, want %q", got, want)
	}
}

func TestNewV3(t *testing.T) {
	name := []byte("example.com")
	u := uuid.NewV3(uuid.DNS, name)
	if got := u[6] >> 4; got != 3 {
		t.Errorf("version = %d, want 3", got)
	}
	if got := u[8] >> 6; got != 0b10 {
		t.Errorf("variant = %02b, want 10", got)
	}
	if got := uuid.NewV3(uuid.DNS, name); got != u {
		t.Errorf("same namespace and name produced %v, want %v", got, u)
	}
	if got := uuid.NewV3(uuid.URL, name); got == u {
		t.Errorf("different namespaces produced the same UUID %v", got)
	}
	if got := uuid.NewV3(uuid.DNS, []byte("example.net")); got == u {
		t.Errorf("different names produced the same UUID %v", got)
	}
}

func TestNewV3InvalidNamespace(t *testing.T) {
	if got := uuid.NewV3(uuid.Namespace(0), []byte("example.com")); got != uuid.Nil() {
		t.Fatalf("NewV3() = %v for an unsupported namespace, want Nil UUID", got)
	}
}
