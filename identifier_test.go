package identifier_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/tozd/go/errors"

	"gitlab.com/tozd/identifier"
)

// TODO: Convert to a fuzzing test and a benchmark.

func TestFromUUID(t *testing.T) {
	t.Parallel()

	for i := 0; i < 100000; i++ {
		u := uuid.New()
		i := identifier.UUID(u)
		assert.Len(t, i, 16)
		s := i.String()
		assert.Len(t, s, 22)
		require.True(t, identifier.Valid(s))
		ii, errE := identifier.MaybeString(s)
		require.NoError(t, errE, "% -+#.1v", errE)
		assert.Equal(t, i, ii)
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	for i := 0; i < 100000; i++ {
		i := identifier.New()
		assert.Len(t, i, 16)
		s := i.String()
		assert.Len(t, s, 22)
		require.True(t, identifier.Valid(s))
		ii, errE := identifier.MaybeString(s)
		require.NoError(t, errE, "% -+#.1v", errE)
		assert.Equal(t, i, ii)
	}
}

func TestValid(t *testing.T) {
	t.Parallel()

	assert.False(t, identifier.Valid(""))
	assert.False(t, identifier.Valid("42"))
	assert.True(t, identifier.Valid("CDEFGHJKLMNPQRSTUVWXYZ"))
	assert.False(t, identifier.Valid("zzzzzzzzzzzzzzzzzzzzzz"))
	assert.True(t, identifier.Valid("2222222222222222222222"))
	assert.True(t, identifier.Valid("2111111111111111111111"))
	assert.True(t, identifier.Valid("1111111111111111211111"))
	assert.True(t, identifier.Valid("1111111111111111111111"))
}

type testStruct struct {
	ID identifier.Identifier `json:"id"`
}

func TestJSON(t *testing.T) {
	t.Parallel()

	x := testStruct{
		ID: identifier.New(),
	}
	data, err := json.Marshal(x)
	require.NoError(t, err)
	assert.Equal(t, `{"id":"`+x.ID.String()+`"}`, string(data))
	var y testStruct
	err = json.Unmarshal(data, &y)
	require.NoError(t, err)
	assert.Equal(t, x.ID, y.ID)
}

func TestFromStringError(t *testing.T) {
	t.Parallel()

	_, errE := identifier.MaybeString("xxx")
	require.Error(t, errE, "% -+#.1v", errE)
	assert.Equal(t, "xxx", errors.AllDetails(errE)["value"])

	assert.PanicsWithError(t, identifier.ErrInvalidIdentifier.Error(), func() {
		identifier.String("xxx")
	})
}

func TestGoStringer(t *testing.T) {
	t.Parallel()

	i := identifier.String("Xuw7QMx5Qqee5jn6ddXCrc")
	s := fmt.Sprintf("%#v", i)
	assert.Equal(t, `identifier.String("Xuw7QMx5Qqee5jn6ddXCrc")`, s)
}

func TestFrom(t *testing.T) {
	t.Parallel()

	t.Run("produces valid identifier", func(t *testing.T) {
		t.Parallel()

		i := identifier.From("test")
		assert.Len(t, i, 16)
		s := i.String()
		assert.Len(t, s, 22)
		require.True(t, identifier.Valid(s))
		assert.Equal(t, "LhYXZThRsu1RXG5ddRZiUt", s)
	})

	t.Run("deterministic", func(t *testing.T) {
		t.Parallel()

		i1 := identifier.From("value1", "value2", "value3")
		i2 := identifier.From("value1", "value2", "value3")
		assert.Equal(t, i1, i2)
		assert.Equal(t, "UhsHqGT45sDirscEZLxmC3", i1.String())
	})

	t.Run("different inputs produce different outputs", func(t *testing.T) {
		t.Parallel()

		i1 := identifier.From("value1")
		i2 := identifier.From("value2")
		assert.NotEqual(t, i1, i2)
		assert.Equal(t, "8UwJ16f3LZEDo1EWEPR1Ua", i1.String())
		assert.Equal(t, "1eNbijZLjE6RCP9J3v6yz1", i2.String())
	})

	t.Run("order matters", func(t *testing.T) {
		t.Parallel()

		i1 := identifier.From("value1", "value2")
		i2 := identifier.From("value2", "value1")
		assert.NotEqual(t, i1, i2)
		assert.Equal(t, "ReKxivb3BXqpCurBhx657A", i1.String())
		assert.Equal(t, "FtRyFysugNvJjVkMztmAuW", i2.String())
	})

	t.Run("single vs multiple values", func(t *testing.T) {
		t.Parallel()

		i1 := identifier.From("value1value2")
		i2 := identifier.From("value1", "value2")
		assert.NotEqual(t, i1, i2)
		assert.Equal(t, "21DW9wo4kBwGXVxbPW69oQ", i1.String())
		assert.Equal(t, "ReKxivb3BXqpCurBhx657A", i2.String())
	})

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()

		i := identifier.From("")
		assert.Len(t, i, 16)
		s := i.String()
		require.True(t, identifier.Valid(s))
		assert.Equal(t, "V7jseQevszwMPhi4evidTR", s)
	})

	t.Run("multiple empty strings", func(t *testing.T) {
		t.Parallel()

		i1 := identifier.From("", "")
		i2 := identifier.From("")
		assert.NotEqual(t, i1, i2)
		assert.Equal(t, "Cbyu7w2KmnA6ZJVbgsHpHH", i1.String())
		assert.Equal(t, "V7jseQevszwMPhi4evidTR", i2.String())
	})

	t.Run("unicode normalization NFC", func(t *testing.T) {
		t.Parallel()

		// U+00E9 (é) vs U+0065 U+0301 (e + combining acute accent)
		// Both should normalize to U+00E9 in NFC.
		i1 := identifier.From("\u00e9")       // é as single character.
		i2 := identifier.From("\u0065\u0301") // e + combining acute.
		assert.Equal(t, i1, i2, "NFC normalization should make these equal")
		assert.Equal(t, "ADHVGvUx5PLGDsjwjY5BA9", i1.String())
	})

	t.Run("specific known values", func(t *testing.T) {
		t.Parallel()

		// Test with known values to ensure consistency across implementations.
		i := identifier.From("test", "value")
		s := i.String()
		// This is just to verify the implementation produces consistent results.
		assert.Len(t, s, 22)
		require.True(t, identifier.Valid(s))
		assert.Equal(t, "J1oVAcLajL9m5GgBJ1eeqz", s)

		// Same input should always produce same output.
		i2 := identifier.From("test", "value")
		assert.Equal(t, i, i2)
	})

	t.Run("cascading hash", func(t *testing.T) {
		t.Parallel()

		// Each additional value should change the hash.
		i1 := identifier.From("a")
		i2 := identifier.From("a", "b")
		i3 := identifier.From("a", "b", "c")
		assert.NotEqual(t, i1, i2)
		assert.NotEqual(t, i2, i3)
		assert.NotEqual(t, i1, i3)
		assert.Equal(t, "S1yrYnjHbfbiTySsN9h1eC", i1.String())
		assert.Equal(t, "KorZ8VDpKQvHrZd2njXraU", i2.String())
		assert.Equal(t, "469q6wDNXV222gSefVXCrJ", i3.String())
	})
}
