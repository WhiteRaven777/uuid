package uuid

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"reflect"
	"strings"
	"time"
)

// UUID is a 128bit Universal Unique IDentifier.
type UUID [16]byte

func (uuid UUID) Version() (ret Version) {
	return Version((binary.BigEndian.Uint64(uuid[:8]) & (1<<16 - 1)) >> 12)
}

func (uuid UUID) Short() (ret string) {
	return hex.EncodeToString(uuid[:])
}

func (uuid UUID) String() (ret string) {
	var buf [36]byte
	hex.Encode(buf[:8], uuid[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], uuid[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], uuid[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], uuid[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], uuid[10:])
	return string(buf[:])
}

func (uuid UUID) URN() (ret string) {
	return "urn:uuid:" + uuid.String()
}

func (uuid UUID) Base64() (ret string) {
	return base64.RawStdEncoding.EncodeToString(uuid[:])
}

func byteDecode(id []byte) (uuid UUID, err error) {
	if id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		return UUID{}, errors.New("invalid UUID format")
	}
	var src [32]byte
	copy(src[:8], id[:8])
	copy(src[8:12], id[9:13])
	copy(src[12:16], id[14:18])
	copy(src[16:20], id[19:23])
	copy(src[20:], id[24:36])
	_, err = hex.Decode(uuid[:], src[:])
	return
}

func Parse(in interface{}) (uuid UUID, err error) {
	if id, ok := in.(string); ok {
		switch len(id) {
		case 36 - 4:
			// short format
			// xxxxxxxxxxxxxxxxxxNxxxxxxxxxxxxxxx
			_, err = hex.Decode(uuid[:], []byte(id))
		case 36:
			// UUID format
			// xxxxxxxxxx-xxxx-xxxx-Nxxx-xxxxxxxxxxxx
			return byteDecode([]byte(id))
		case 36 + 2:
			//  // Microsoft format -> undefined
			//	// {xxxxxxxxxx-xxxx-xxxx-Nxxx-xxxxxxxxxxxx}
			if id[0] != '{' || id[37] != '}' {
				return UUID{}, errors.New("invalid UUID format")
			}
			return byteDecode([]byte(id[1:37]))
		case 36 + 9:
			// URN format
			// urn:uuid:xxxxxxxxxx-xxxx-xxxx-Nxxx-xxxxxxxxxxxx
			if strings.ToLower(id[:9]) != "urn:uuid:" {
				return UUID{}, errors.New("invalid UUID format")
			}
			return byteDecode([]byte(id[9:]))
		default:
			// base64 ?
			if _, err = base64.RawStdEncoding.Decode(uuid[:], []byte(id)); err != nil {
				if _, err = base64.StdEncoding.Decode(uuid[:], []byte(id)); err != nil {
					return
				}
			}
		}
	} else {
		if x, ok := in.([]uint8); ok {
			copy(uuid[:], x)
		} else {
			err = errors.New("input value is an unsupported type")
		}
	}
	return
}

// Version is used to represent the version of the UUID.
type Version byte

const (
	UUIDv1 = Version(0b0001)
	UUIDv2 = Version(0b0010)
	UUIDv3 = Version(0b0011)
	UUIDv4 = Version(0b0100)
	UUIDv5 = Version(0b0101)
	UUIDv6 = Version(0b0110)
	UUIDv7 = Version(0b0111)
	UUIDv8 = Version(0b1000)
)

// Variant is used to represent the variant of the UUID.
type Variant byte

const (
	Backward  = Variant(0b_0000)
	RFC4122   = Variant(0b_0010)
	Microsoft = Variant(0b_0110)
	Future    = Variant(0b_0111)
)

var reader = rand.Reader

type ID interface {
	Convert(t time.Time) (ret UUID)
	Generate() (ret UUID)
	Decode(uuid UUID) (ret time.Time)
}

func New(version Version) (id ID) {
	switch version {
	case UUIDv6:
	case UUIDv7:
		id = &V7{v: version}
	case UUIDv8:
	}
	return
}

func getMonotonicClock(t time.Time, max ...int) (ret uint64) {
	const hasMonotonic = 1 << 63
	var (
		wall = reflect.ValueOf(t).FieldByName("wall").Uint()
		ext  = reflect.ValueOf(t).FieldByName("ext").Int()
	)
	if wall&hasMonotonic == 0 {
		return 0
	}

	ret = uint64(ext)
	if len(max) > 0 {
		ret = ret & (1<<max[0] - 1)
	}
	return
}
