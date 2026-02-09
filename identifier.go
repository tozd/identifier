// Package identifier provides functions to generate and parse readable global identifiers.
package identifier

import (
	"crypto/rand"
	"crypto/sha256"
	"io"
	"regexp"
	"strings"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/google/uuid"
	"gitlab.com/tozd/go/errors"
	"golang.org/x/text/unicode/norm"
)

const (
	stringLength   = 22
	bytesMinLength = 16
)

var ErrInvalidIdentifier = errors.Base("invalid identifier")

var idRegex = regexp.MustCompile(`^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]{22}$`)

// Identifier is a 128-bit identifier.
//
//nolint:recvcheck
type Identifier [16]byte

// String encodes Identifier value into a string using base 58 encoding.
func (i Identifier) String() string {
	res := base58.Encode(i[:])
	if len(res) < stringLength {
		// String might be shorter than stringLength to encode 128 bits, in that
		// we do zero left padding (character "1" in base 58).
		return strings.Repeat("1", stringLength-len(res)) + res
	}
	return res
}

// UnmarshalText implements encoding.TextUnmarshaler interface for Identifier.
func (i *Identifier) UnmarshalText(text []byte) error {
	ii, err := MaybeString(string(text))
	if err != nil {
		return err
	}
	*i = ii
	return nil
}

// MarshalText implements encoding.TextMarshaler interface for Identifier.
func (i Identifier) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// GoString implements fmt.GoStringer interface for Identifier.
func (i Identifier) GoString() string {
	return `identifier.String("` + i.String() + `")`
}

// UUID returns the UUID encoded as an Identifier.
func UUID(data uuid.UUID) Identifier {
	return Data(data)
}

// Data returns 16 bytes data encoded as an Identifier.
func Data(data [16]byte) Identifier {
	return Identifier(data)
}

// MaybeString parses a string-encoded identifier in base 58 encoding
// into a corresponding Identifier value.
func MaybeString(data string) (Identifier, errors.E) {
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
	return Identifier([16]byte(res[len(res)-16:])), nil
}

// String is the same as MaybeString but panics on an error.
func String(data string) Identifier {
	i, err := MaybeString(data)
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

// Valid returns true if id string is a valid identifier
// (MaybeString will not return an error).
func Valid(id string) bool {
	if !idRegex.MatchString(id) {
		return false
	}
	_, err := MaybeString(id)
	return err == nil
}

// From generates a deterministic identifier from one or more string values using iterative SHA-256 hashing.
//
// Each value is normalized using Unicode NFC normalization before hashing. The function computes
// hash = SHA256(normalize(values[0])), then hash = SHA256(hash + normalize(values[1])), and so on.
// The final identifier is derived from the first 128 bits of the resulting hash.
//
// Different values or different orderings produce different identifiers.
func From(values ...string) Identifier {
	var hash []byte
	for i, value := range values {
		// Normalize the string using NFC.
		normalized := norm.NFC.String(value)

		if i == 0 {
			// First iteration: hash just the normalized value.
			h := sha256.Sum256([]byte(normalized))
			hash = h[:]
		} else {
			// Subsequent iterations: hash = sha256(hash + normalized).
			hash = append(hash, []byte(normalized)...)
			h := sha256.Sum256(hash)
			hash = h[:]
		}
	}

	// Take first 128 bits (16 bytes) of the final hash.
	var result [16]byte
	copy(result[:], hash[:16])
	return Identifier(result)
}
