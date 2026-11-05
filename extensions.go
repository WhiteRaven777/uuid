package uuid

import (
	"encoding/base64"
	"encoding/hex"
	"time"
)

// Short returns the UUID as 32 lowercase hexadecimal digits without hyphens.
func (u UUID) Short() string { return hex.EncodeToString(u[:]) }

// URN returns the UUID in its URN representation.
func (u UUID) URN() string { return "urn:uuid:" + u.String() }

// Base64 returns the UUID using unpadded standard base64 encoding.
func (u UUID) Base64() string { return base64.RawStdEncoding.EncodeToString(u[:]) }

// Time returns the time embedded in a version 1, 2, 6, or 7 UUID.
//
// For UUIDv2, Time returns the beginning of the approximately seven-minute
// interval represented by the UUID because its low 32 timestamp bits are
// replaced by the local identifier. For UUID versions without a defined time
// field, Time returns the zero time and false.
func (u UUID) Time() (t time.Time, ok bool) {
	var timestamp uint64
	switch u[6] >> 4 {
	case 1:
		timestamp = uint64(u[0])<<24 | uint64(u[1])<<16 | uint64(u[2])<<8 | uint64(u[3])
		timestamp |= uint64(u[4])<<40 | uint64(u[5])<<32
		timestamp |= uint64(u[6]&0x0f)<<56 | uint64(u[7])<<48
		t = gregorianTime(timestamp)
	case 2:
		timestamp = uint64(u[4])<<40 | uint64(u[5])<<32
		timestamp |= uint64(u[6]&0x0f)<<56 | uint64(u[7])<<48
		t = gregorianTime(timestamp)
	case 6:
		timestamp = uint64(u[0])<<52 | uint64(u[1])<<44 | uint64(u[2])<<36 | uint64(u[3])<<28
		timestamp |= uint64(u[4])<<20 | uint64(u[5])<<12
		timestamp |= uint64(u[6]&0x0f)<<8 | uint64(u[7])
		t = gregorianTime(timestamp)
	case 7:
		timestamp = uint64(u[0])<<40 | uint64(u[1])<<32 | uint64(u[2])<<24
		timestamp |= uint64(u[3])<<16 | uint64(u[4])<<8 | uint64(u[5])
		t = time.UnixMilli(int64(timestamp)).UTC()
	default:
		return
	}
	ok = true
	return
}

// gregorianTime converts UUID timestamp ticks to UTC.
func gregorianTime(timestamp uint64) time.Time {
	ticks := int64(timestamp) - int64(gregorianEpochOffset)
	return time.Unix(ticks/10_000_000, ticks%10_000_000*100).UTC()
}
