// Package provides functions to generate and parse readable global identifiers.
package identifier

import (
	"crypto/rand"
	"io"
	"regexp"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"gitlab.com/tozd/go/errors"
)

const (
	stringLength   = 22
	bytesMinLength = 16
)

var ErrInvalidIdentifier = errors.Base("invalid identifier")

var idRegex = regexp.MustCompile(`^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]{22}$`)

type Identifier [16]byte

// String encodes Identifier value into a string using base 58 encoding.
func (i Identifier) String() string {
	res := base58.Encode(i[:])
	if len(res) < stringLength {
		// String might be shorter than stringLength to encode 128 bits, in that
		// we do zero left padding (character "1" in base58).
		return strings.Repeat("1", stringLength-len(res)) + res
	}
	return res
}

func (i *Identifier) UnmarshalText(text []byte) error {
	ii, err := FromString(string(text))
	if err != nil {
		return err
	}
	*i = ii
	return nil
}

func (i Identifier) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// FromUUID returns the UUID encoded as an identifier.
func FromUUID(data uuid.UUID) Identifier {
	return FromData(data)
}

// FromData returns 16 bytes data encoded as an identifier.
func FromData(data [16]byte) Identifier {
	return Identifier(data)
}

// FromString parses a string-encoded identifier in base 58 encoding
// into a corresponding Identifier value.
func FromString(data string) (Identifier, errors.E) {
	if len(data) != stringLength {
		return Identifier{}, errors.WithDetails(ErrInvalidIdentifier, "value", data)
	}
	res := base58.Decode(data)
	// Decode returns an empty slice if data contains a character outside of base58.
	// But we care about too short strings here too.
	if len(res) < bytesMinLength {
		return Identifier{}, errors.WithDetails(ErrInvalidIdentifier, "value", data)
	}
	// String might longer than necessary to encode 128 bits, in that case we require extra bytes
	// at the beginning to be zero (or character "1" in base58), i.e., zero left padding.
	for i := 0; i+bytesMinLength < len(res); i++ {
		if res[i] != 0 {
			return Identifier{}, errors.WithDetails(ErrInvalidIdentifier, "value", data)
		}
	}
	// We take the last 16 bytes.
	return Identifier(*(*[16]byte)(res[len(res)-16:])), nil
}

// MustFromString is the same as FromString but panics on an error.
func MustFromString(data string) Identifier {
	i, err := FromString(data)
	if err != nil {
		panic(err)
	}
	return i
}

// New returns a new random identifier.
func New() Identifier {
	return MustFromReader(rand.Reader)
}

// FromReader returns a new random identifier using r as a source of randomness.
func FromReader(r io.Reader) (Identifier, errors.E) {
	// We read 128 bits.
	data := [16]byte{}
	_, err := io.ReadFull(r, data[:])
	if err != nil {
		return Identifier{}, errors.WithStack(err)
	}
	return Identifier(data), nil
}

// MustFromReader is the same as FromReader but panics on an error.
func MustFromReader(r io.Reader) Identifier {
	i, err := FromReader(r)
	if err != nil {
		panic(err)
	}
	return i
}

// Valid returns true if id string is a valid identifier.
func Valid(id string) bool {
	if !idRegex.MatchString(id) {
		return false
	}
	_, err := FromString(id)
	return err == nil
}
