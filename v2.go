package uuid

import (
	"encoding/binary"
	"sync"
	"time"
)

var (
	randomV2Generator dceGenerator
	macV2Generator    dceGenerator
)

type dceGenerator struct {
	mu sync.Mutex

	gregorian   gregorianGenerator
	initialised bool
	sequence    uint8
}

// LocalDomain identifies the meaning of a UUIDv2 local identifier.
type LocalDomain uint8

const (
	// Person identifies a POSIX user ID.
	Person LocalDomain = iota
	// Group identifies a POSIX group ID.
	Group
	// Organisation identifies an organisation-defined local identifier.
	Organisation
)

// NewV2 returns a new DCE Security version 2 UUID.
//
// localID is interpreted according to domain. Person and Group identify POSIX
// user and group IDs respectively, while Organisation identifies a value whose
// meaning is defined by the organisation using it.
//
// UUIDv2 replaces the low 32 bits of the UUIDv1 timestamp with localID and the
// low 8 bits of its clock sequence with domain. Consequently, its timestamp
// has a resolution of approximately seven minutes and only 6 clock-sequence
// bits remain. Callers must account for this limited uniqueness.
//
// By default, the node field is randomly generated so that UUIDs do not
// disclose the host's MAC address. If useMAC is true, NewV2 uses a MAC address
// from a network interface when one is available and otherwise falls back to a
// randomly generated node identifier.
//
// NewV2 returns the Nil UUID if domain is not one of the LocalDomain constants
// defined by this package.
func NewV2(domain LocalDomain, localID uint32, useMAC ...bool) UUID {
	if domain > Organisation {
		return Nil()
	}
	if len(useMAC) > 0 && useMAC[0] {
		return formatV2(formatV1(macV2Generator.next(time.Now(), hardwareNode)), domain, localID)
	}
	return formatV2(formatV1(randomV2Generator.next(time.Now(), nil)), domain, localID)
}

// next advances the 6 clock-sequence bits retained by the UUIDv2 layout.
func (g *dceGenerator) next(now time.Time, source nodeSource) (timestamp uint64, clockSequence uint16, node [6]byte) {
	g.mu.Lock()
	defer g.mu.Unlock()

	timestamp, clockSequence, node = g.gregorian.next(now, source)
	if !g.initialised {
		g.sequence = uint8(clockSequence >> 8)
		g.initialised = true
	}
	clockSequence = uint16(g.sequence) << 8
	g.sequence = (g.sequence + 1) & 0x3f
	return
}

// formatV2 replaces UUIDv1 fields with a DCE local identifier and domain.
func formatV2(u UUID, domain LocalDomain, localID uint32) UUID {
	binary.BigEndian.PutUint32(u[0:4], localID)
	u[6] = u[6]&0x0f | 0x20
	u[9] = byte(domain)
	return u
}
