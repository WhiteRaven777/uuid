package uuid

import (
	"encoding/binary"
	"time"
)

var (
	randomV6Generator gregorianGenerator
	macV6Generator    gregorianGenerator
)

// NewV6 returns a new version 6 UUID.
//
// UUIDv6 uses the same Gregorian timestamp, clock sequence, and node fields as
// UUIDv1, but orders the timestamp from most to least significant for improved
// database locality.
//
// By default, the node field is randomly generated so that UUIDs do not
// disclose the host's MAC address. If useMAC is true, NewV6 uses a MAC address
// from a network interface when one is available and otherwise falls back to a
// randomly generated node identifier.
func NewV6(useMAC ...bool) UUID {
	var timestamp uint64
	var clockSequence uint16
	var node [6]byte
	if len(useMAC) > 0 && useMAC[0] {
		timestamp, clockSequence, node = macV6Generator.next(time.Now(), hardwareNode)
	} else {
		timestamp, clockSequence, node = randomV6Generator.next(time.Now(), nil)
	}
	return formatV6(timestamp, clockSequence, node)
}

func formatV6(timestamp uint64, clockSequence uint16, node [6]byte) UUID {
	var u UUID
	binary.BigEndian.PutUint32(u[0:4], uint32(timestamp>>28))
	binary.BigEndian.PutUint16(u[4:6], uint16(timestamp>>12))
	binary.BigEndian.PutUint16(u[6:8], uint16(timestamp)&0x0fff|0x6000)
	binary.BigEndian.PutUint16(u[8:10], clockSequence&0x3fff|0x8000)
	copy(u[10:], node[:])
	return u
}
