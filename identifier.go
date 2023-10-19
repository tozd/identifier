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

// FromUUID returns an UUID encoded as an identifier.
func FromUUID(data uuid.UUID) Identifier {
	return Identifier(data)
}

func FromString(data string) Identifier {
	res := base58.Decode(data)
	if len(res) < 16 {
		panic(errors.Errorf(`invalid identifier data length: %d`, len(res)))
	}
	for i := 0; i+16 < len(res); i++ {
		if res[i] != 0 {
			panic(errors.Errorf(`invalid extra byte: %x`, res[i]))
		}
	}
	// We take the last 16 bytes.
	return Identifier(*(*[16]byte)(res[len(res)-16:]))
}

// New returns a new random identifier.
func New() Identifier {
	return FromReader(rand.Reader)
}

// NewRandom returns a new random identifier using r as a source of randomness.
func FromReader(r io.Reader) Identifier {
	// We read 128 bits.
	data := [16]byte{}
	_, err := io.ReadFull(r, data[:])
	if err != nil {
		panic(errors.WithStack(err))
	}
	return Identifier(data)
}

// Valid returns true if id string looks like a valid identifier.
func Valid(id string) (res bool) { //nolint:nonamedreturns
	if !idRegex.MatchString(id) {
		return false
	}
	defer func() {
		if recover() != nil {
			res = false
		}
	}()
	FromString(id)
	return true
}
