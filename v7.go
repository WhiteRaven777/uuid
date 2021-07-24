package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"math/big"
	"sync"
	"time"
)

// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                            unixts                             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |unixts |       subsec_a        |  ver  |       subsec_b        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |var|                   subsec_seq_node                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                       subsec_seq_node                         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

type Precision int

const (
	PrecisionDefault = Precision(iota)
	// => High

	PrecisionLow
	// Low: msec precision
	// 0                   1                   2                   3
	// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                            unixts                             |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |unixts |         msec          |  ver  |          seq          |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |var|                         rand                              |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                             rand                              |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

	PrecisionMid
	// Low: usec precision
	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                            unixts                             |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |unixts |         usec          |  ver  |         usec          |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |var|             seq           |            rand               |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                             rand                              |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

	PrecisionHigh
	// Low: nsec precision
	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                            unixts                             |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |unixts |         nsec          |  ver  |         nsec          |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |var|             nsec          |      seq      |     rand      |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |                             rand                              |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
)

func (p Precision) String() (ret string) {
	switch p {
	case PrecisionDefault:
		ret = "high (default)"
	case PrecisionLow:
		ret = "low"
	case PrecisionMid:
		ret = "middle"
	case PrecisionHigh:
		ret = "high"
	default:
		ret = "undefine"
	}
	return
}

type V7 struct {
	v  Version
	p  Precision
	mx sync.Mutex
}

const (
	v7UnixBit    = 36
	v7SubSecA    = 12
	v7SubSecB    = 12
	v7SubSecC    = 14
	v7VersionBit = 4
	v7VariantBit = 2
	v7LowSeq     = 12
	v7MidSeq     = 14
	v7HighSeq    = 8
	v7LowRand    = 32*2 - 2
	v7MidRand    = 16 + 32
	v7HighRand   = 8 + 32
)

func (id *V7) SetPrecision(p Precision) {
	id.mx.Lock()
	defer id.mx.Unlock()
	id.p = p
	return
}

func (id *V7) rand() (r uint64) {
	var max int64
	switch id.p {
	case PrecisionLow:
		// msec precision
		max = 1<<v7LowRand - 1
	case PrecisionMid:
		// usec precision
		max = 1<<v7MidRand - 1
	case PrecisionHigh, PrecisionDefault:
		// nsec precision
		max = 1<<v7HighRand - 1
	}

	if n, err := rand.Int(reader, big.NewInt(max)); err == nil {
		r = n.Uint64()
	}
	return
}

func (id *V7) convert(now time.Time) (ret UUID) {
	un := uint64(now.UnixNano())

	var (
		up   uint64
		down uint64
		r    = id.rand()
	)

	switch id.p {
	case PrecisionLow:
		// msec precision
		var (
			uxts = un / 1000 / 1000 / 1000
			msec = (un / 1000 / 1000) % 1000
		)

		up = (uxts << (v7SubSecA + v7VersionBit + v7SubSecB)) |
			(msec << (v7VersionBit + v7SubSecB)) |
			(uint64(id.v) << (v7SubSecB)) |
			(getMonotonicClock(now, v7LowSeq))

		down = (uint64(RFC4122) << (v7LowRand)) |
			(r)
	case PrecisionMid:
		// usec precision
		var (
			uxts = un / 1000 / 1000 / 1000
			msec = (un / 1000 / 1000) % 1000
			usec = (un / 1000) % 1000
		)

		up = (uxts << (v7SubSecA + v7VersionBit + v7SubSecB)) |
			(msec << (v7VersionBit + v7SubSecB)) |
			(uint64(id.v) << (v7SubSecB)) |
			(usec)

		down = (uint64(RFC4122) << (v7MidSeq + v7MidRand)) |
			(getMonotonicClock(now, v7MidSeq) << (v7MidRand)) |
			(r)
	case PrecisionHigh, PrecisionDefault:
		// nsec precision
		var (
			uxts = un / 1000 / 1000 / 1000
			msec = (un / 1000 / 1000) % 1000
			usec = (un / 1000) % 1000
			nsec = un % 1000
		)

		up = (uxts << (v7SubSecA + v7VersionBit + v7SubSecB)) |
			(msec << (v7VersionBit + v7SubSecB)) |
			(uint64(id.v) << (v7SubSecB)) |
			(usec)

		down = (uint64(RFC4122) << (v7SubSecC + v7HighSeq + v7HighRand)) |
			(nsec << (v7HighSeq + v7HighRand)) |
			(getMonotonicClock(now, v7HighSeq) << (v7HighRand)) |
			(r)
	}

	binary.BigEndian.PutUint64(ret[:8], up)
	binary.BigEndian.PutUint64(ret[8:], down)
	return
}

func (id *V7) Convert(t time.Time) (ret UUID) {
	return id.convert(t)
}

func (id *V7) Generate() (ret UUID) {
	return id.convert(time.Now())
}

func (id *V7) Decode(uuid UUID) (ret time.Time) {
	var (
		version = (binary.BigEndian.Uint64(uuid[:8]) & (1<<(v7VersionBit+v7SubSecB) - 1)) >> v7SubSecB
		variant = binary.BigEndian.Uint64(uuid[8:]) >> (64 - v7VariantBit)

		uxts = binary.BigEndian.Uint64(uuid[:8]) >> (64 - v7UnixBit)
		msec = (binary.BigEndian.Uint64(uuid[:8]) & (1<<(v7SubSecA+v7VersionBit+v7SubSecB) - 1)) >> (v7VersionBit + v7SubSecB)
		usec = binary.BigEndian.Uint64(uuid[:8]) & (1<<v7SubSecB - 1)
		nsec = binary.BigEndian.Uint64(uuid[8:]) & (1<<(v7SubSecC+v7HighSeq+v7HighRand) - 1) >> (v7HighSeq + v7HighRand)
	)

	if Version(version) != UUIDv7 {
		return
	}

	if Variant(variant) != RFC4122 {
		return
	}

	if msec >= 1000 {
		msec = 0
	}
	if usec >= 1000 {
		usec = 0
	}
	if nsec >= 1000 {
		nsec = 0
	}

	return time.Unix(int64(uxts), int64(msec*1000*1000+usec*1000+nsec)).UTC()
}
