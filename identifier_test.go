package identifier_test

import (
	"encoding/json"
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
		i := identifier.FromUUID(u)
		assert.Len(t, i, 16)
		s := i.String()
		assert.Len(t, s, 22)
		require.True(t, identifier.Valid(s))
		ii, errE := identifier.FromString(s)
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
		ii, errE := identifier.FromString(s)
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

	_, errE := identifier.FromString("xxx")
	require.Error(t, errE, "% -+#.1v", errE)
	assert.Equal(t, "xxx", errors.AllDetails(errE)["value"])

	assert.PanicsWithError(t, identifier.ErrInvalidIdentifier.Error(), func() {
		identifier.MustFromString("xxx")
	})
}
