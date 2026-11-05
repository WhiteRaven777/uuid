package uuid

import "encoding/binary"

const (
	customAMask = uint64(1<<48 - 1)
	customBMask = uint16(1<<12 - 1)
	customCMask = uint64(1<<62 - 1)
)

// NewV8 returns a custom version 8 UUID.
//
// UUIDv8 leaves 122 bits to an application-specific format. The valid custom
// field ranges are 0x000000000000 to 0xffffffffffff for customA, 0x000 to
// 0xfff for customB, and 0x0000000000000000 to 0x3fffffffffffffff for customC.
// Values outside these ranges are masked to the field width, discarding their
// more significant bits. The version and variant bits are set by this
// function.
//
// UUIDv8 does not define or guarantee uniqueness. Callers are responsible for
// defining the custom fields and their uniqueness properties. Use NewV4
// instead when all 122 custom bits would contain only random data.
func NewV8(customA uint64, customB uint16, customC uint64) (u UUID) {
	binary.BigEndian.PutUint64(u[0:8], customA&customAMask<<16|0x8000|uint64(customB&customBMask))
	binary.BigEndian.PutUint64(u[8:16], customC&customCMask|0x8000000000000000)
	return
}
