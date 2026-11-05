package uuid

import "crypto/md5"

// Namespace identifies a standard namespace for name-based UUIDs.
type Namespace uint8

const (
	// DNS identifies the Domain Name System namespace.
	DNS Namespace = iota + 1
	// URL identifies the Uniform Resource Locator namespace.
	URL
	// OID identifies the ISO object identifier namespace.
	OID
	// X500 identifies the X.500 Distinguished Name namespace.
	X500
)

// NewV3 returns a version 3 UUID derived from namespace and name using MD5.
//
// NewV3 does not canonicalise name. Callers must provide the canonical byte
// representation defined by the selected namespace.
//
// UUIDv3 uses MD5 and is intended for compatibility, not security-sensitive
// applications.
//
// NewV3 returns the Nil UUID if namespace is not one of the Namespace constants
// defined by this package.
func NewV3(namespace Namespace, name []byte) (u UUID) {
	if u = namespaceUUID(namespace); u == Nil() {
		return
	}
	hash := md5.New()
	hash.Write(u[:])
	hash.Write(name)

	copy(u[:], hash.Sum(nil))
	u[6] = u[6]&0x0f | 0x30
	u[8] = u[8]&0x3f | 0x80
	return
}

// namespaceUUID resolves a public namespace mode to its RFC-defined UUID.
func namespaceUUID(namespace Namespace) (u UUID) {
	switch namespace {
	case DNS:
		// DNS namespace: 6ba7b810-9dad-11d1-80b4-00c04fd430c8
		u = UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	case URL:
		// URL namespace: 6ba7b811-9dad-11d1-80b4-00c04fd430c8
		u = UUID{0x6b, 0xa7, 0xb8, 0x11, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	case OID:
		// OID namespace: 6ba7b812-9dad-11d1-80b4-00c04fd430c8
		u = UUID{0x6b, 0xa7, 0xb8, 0x12, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	case X500:
		// X.500 DN namespace: 6ba7b814-9dad-11d1-80b4-00c04fd430c8
		u = UUID{0x6b, 0xa7, 0xb8, 0x14, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	}
	return
}
