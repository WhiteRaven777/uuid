package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"net"
	"sync"
	"time"
)

const gregorianEpochOffset = uint64(122192928000000000)

var (
	randomV1Generator gregorianGenerator
	macV1Generator    gregorianGenerator
)

type nodeSource func() (node [6]byte, ok bool)

type gregorianGenerator struct {
	mu sync.Mutex

	initialised   bool
	node          [6]byte
	clockSequence uint16
	lastWall      uint64
	lastTimestamp uint64
}

// NewV1 returns a new version 1 UUID.
//
// By default, the node field is randomly generated so that UUIDs do not
// disclose the host's MAC address. If useMAC is true, NewV1 uses a MAC address
// from a network interface when one is available and otherwise falls back to a
// randomly generated node identifier.
func NewV1(useMAC ...bool) UUID {
	var timestamp uint64
	var clockSequence uint16
	var node [6]byte
	if len(useMAC) > 0 && useMAC[0] {
		timestamp, clockSequence, node = macV1Generator.next(time.Now(), hardwareNode)
	} else {
		timestamp, clockSequence, node = randomV1Generator.next(time.Now(), nil)
	}
	return formatV1(timestamp, clockSequence, node)
}

func (g *gregorianGenerator) next(now time.Time, source nodeSource) (timestamp uint64, clockSequence uint16, node [6]byte) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.initialised {
		g.initialise(source)
	}

	wall := gregorianTimestamp(now)
	if g.lastWall != 0 && wall < g.lastWall {
		g.clockSequence = (g.clockSequence + 1) & 0x3fff
	}

	timestamp = wall
	if timestamp <= g.lastTimestamp {
		timestamp = g.lastTimestamp + 1
	}

	g.lastWall = wall
	g.lastTimestamp = timestamp
	return timestamp, g.clockSequence, g.node
}

func (g *gregorianGenerator) initialise(source nodeSource) {
	if source != nil {
		if node, ok := source(); ok {
			g.node = node
			g.initialiseClockSequence()
			g.initialised = true
			return
		}
	}

	var seed [8]byte
	rand.Read(seed[:])

	copy(g.node[:], seed[:6])
	// Set the multicast bit to mark the node as randomly generated.
	g.node[0] |= 0x01
	g.clockSequence = binary.BigEndian.Uint16(seed[6:]) & 0x3fff
	g.initialised = true
}

func (g *gregorianGenerator) initialiseClockSequence() {
	var seed [2]byte
	rand.Read(seed[:])
	g.clockSequence = binary.BigEndian.Uint16(seed[:]) & 0x3fff
}

func hardwareNode() (node [6]byte, ok bool) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return node, false
	}

	for _, networkInterface := range interfaces {
		if networkInterface.Flags&net.FlagUp == 0 || networkInterface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if node, ok = macNode(networkInterface.HardwareAddr); ok {
			return node, true
		}
	}
	for _, networkInterface := range interfaces {
		if node, ok = macNode(networkInterface.HardwareAddr); ok {
			return node, true
		}
	}
	return node, false
}

func macNode(address net.HardwareAddr) (node [6]byte, ok bool) {
	if len(address) != len(node) || address[0]&0x01 != 0 {
		return node, false
	}
	copy(node[:], address)
	if node == [6]byte{} {
		return [6]byte{}, false
	}
	return node, true
}

func gregorianTimestamp(t time.Time) uint64 {
	return gregorianEpochOffset + uint64(t.Unix())*10_000_000 + uint64(t.Nanosecond()/100)
}

func formatV1(timestamp uint64, clockSequence uint16, node [6]byte) UUID {
	var u UUID
	binary.BigEndian.PutUint32(u[0:4], uint32(timestamp))
	binary.BigEndian.PutUint16(u[4:6], uint16(timestamp>>32))
	binary.BigEndian.PutUint16(u[6:8], uint16(timestamp>>48)&0x0fff|0x1000)
	binary.BigEndian.PutUint16(u[8:10], clockSequence&0x3fff|0x8000)
	copy(u[10:], node[:])
	return u
}
