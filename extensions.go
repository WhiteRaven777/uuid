package uuid

import (
	"encoding/base64"
	"encoding/hex"
)

// Short returns the UUID as 32 lowercase hexadecimal digits without hyphens.
func (u UUID) Short() string { return hex.EncodeToString(u[:]) }

// URN returns the UUID in its URN representation.
func (u UUID) URN() string { return "urn:uuid:" + u.String() }

// Base64 returns the UUID using unpadded standard base64 encoding.
func (u UUID) Base64() string { return base64.RawStdEncoding.EncodeToString(u[:]) }
