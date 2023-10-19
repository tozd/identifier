// Package provides functions to generate identifiers.
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
	idLength = 22
)

var idRegex = regexp.MustCompile(`^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]{22}$`)

type Identifier [16]byte

func (i Identifier) String() string {
	res := base58.Encode(i[:])
	if len(res) < idLength {
		return strings.Repeat("1", idLength-len(res)) + res
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

// FromUUID returns an UUID encoded as an identifier.
func FromUUID(data uuid.UUID) Identifier {
	return Identifier(data)
}

func FromString(data string) (Identifier, errors.E) {
	res := base58.Decode(data)
	if len(res) < 16 {
		return Identifier{}, errors.Errorf(`invalid identifier data length: %d`, len(res))
	}
	for i := 0; i+16 < len(res); i++ {
		if res[i] != 0 {
			return Identifier{}, errors.Errorf(`invalid extra byte: %x`, res[i])
		}
	}
	// We take the last 16 bytes.
	return Identifier(*(*[16]byte)(res[len(res)-16:])), nil
}

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

// NewRandom returns a new random identifier using r as a source of randomness.
func FromReader(r io.Reader) (Identifier, errors.E) {
	// We read 128 bits.
	data := [16]byte{}
	_, err := io.ReadFull(r, data[:])
	if err != nil {
		return Identifier{}, errors.WithStack(err)
	}
	return Identifier(data), nil
}

func MustFromReader(r io.Reader) Identifier {
	i, err := FromReader(r)
	if err != nil {
		panic(err)
	}
	return i
}

// Valid returns true if id string looks like a valid identifier.
func Valid(id string) bool {
	if !idRegex.MatchString(id) {
		return false
	}
	_, err := FromString(id)
	return err == nil
}
