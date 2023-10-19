package identifier_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/tozd/identifier"
)

// TODO: Convert to a fuzzing test and a benchmark.

func TestFromUUID(t *testing.T) {
	for i := 0; i < 100000; i++ {
		u := uuid.New()
		i := identifier.FromUUID(u)
		assert.Len(t, i, 16)
		s := i.String()
		assert.Len(t, s, 22)
		require.True(t, identifier.Valid(s))
		assert.Equal(t, i, identifier.FromString(s))
	}
}

func TestFromRandom(t *testing.T) {
	for i := 0; i < 100000; i++ {
		i := identifier.New()
		assert.Len(t, i, 16)
		s := i.String()
		assert.Len(t, s, 22)
		require.True(t, identifier.Valid(s))
		assert.Equal(t, i, identifier.FromString(s))
	}
}

func TestValid(t *testing.T) {
	assert.False(t, identifier.Valid(""))
	assert.False(t, identifier.Valid("42"))
	assert.True(t, identifier.Valid("CDEFGHJKLMNPQRSTUVWXYZ"))
	assert.False(t, identifier.Valid("zzzzzzzzzzzzzzzzzzzzzz"))
	assert.True(t, identifier.Valid("2222222222222222222222"))
	assert.True(t, identifier.Valid("2111111111111111111111"))
	assert.True(t, identifier.Valid("1111111111111111211111"))
	assert.True(t, identifier.Valid("1111111111111111111111"))
}
