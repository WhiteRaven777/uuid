package uuid

import "crypto/sha1"

// NewV5 returns a version 5 UUID derived from namespace and name using SHA-1.
//
// NewV5 does not canonicalise name. Callers must provide the canonical byte
// representation defined by the selected namespace.
//
// NewV5 returns the Nil UUID if namespace is not one of the Namespace constants
// defined by this package.
func NewV5(namespace Namespace, name []byte) (u UUID) {
	if u = namespaceUUID(namespace); u == Nil() {
		return
	}
	hash := sha1.New()
	hash.Write(u[:])
	hash.Write(name)

	copy(u[:], hash.Sum(nil)[:len(u)])
	u[6] = u[6]&0x0f | 0x50
	u[8] = u[8]&0x3f | 0x80
	return
}
