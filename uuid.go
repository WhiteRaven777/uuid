// Package uuid extends the Go standard library's uuid package with a few
// additional encodings.
package uuid

import stduuid "uuid"

// UUID is a Universally Unique Identifier as specified in RFC 9562.
// It has the same representation as [uuid.UUID].
type UUID stduuid.UUID

// Parse wraps [uuid.Parse].
func Parse(s string) (UUID, error) {
	u, err := stduuid.Parse(s)
	return UUID(u), err
}

// MustParse wraps [uuid.MustParse].
func MustParse(s string) UUID { return UUID(stduuid.MustParse(s)) }

// New wraps [uuid.New].
func New() UUID { return UUID(stduuid.New()) }

// Nil wraps [uuid.Nil].
func Nil() UUID { return UUID(stduuid.Nil()) }

// Max wraps [uuid.Max].
func Max() UUID { return UUID(stduuid.Max()) }

// NewV4 wraps [uuid.NewV4].
func NewV4() UUID { return UUID(stduuid.NewV4()) }

// NewV7 wraps [uuid.NewV7].
func NewV7() UUID { return UUID(stduuid.NewV7()) }

// String wraps [uuid.UUID.String].
func (u UUID) String() string { return stduuid.UUID(u).String() }

// MarshalText wraps [uuid.UUID.MarshalText].
func (u UUID) MarshalText() ([]byte, error) { return stduuid.UUID(u).MarshalText() }

// AppendText wraps [uuid.UUID.AppendText].
func (u UUID) AppendText(b []byte) ([]byte, error) {
	return stduuid.UUID(u).AppendText(b)
}

// UnmarshalText wraps [uuid.UUID.UnmarshalText].
func (u *UUID) UnmarshalText(b []byte) error {
	var standard stduuid.UUID
	if err := standard.UnmarshalText(b); err != nil {
		return err
	}
	*u = UUID(standard)
	return nil
}

// Compare wraps [uuid.UUID.Compare].
func (u UUID) Compare(v UUID) int {
	return stduuid.UUID(u).Compare(stduuid.UUID(v))
}
