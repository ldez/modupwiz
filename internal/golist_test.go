package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_extractUpdates(t *testing.T) {
	file, err := os.Open(filepath.Join("fixtures", "deps.json"))
	require.NoError(t, err)

	updates, err := extractUpdates(file)
	require.NoError(t, err)

	golden := filepath.Clean(filepath.Join("fixtures", "deps-golden.json"))

	if _, ok := os.LookupEnv("UPDATE_GOLDEN"); ok {
		f, errC := os.Create(golden)
		require.NoError(t, errC)
		t.Cleanup(func() { _ = f.Close() })

		errC = json.NewEncoder(f).Encode(updates)
		require.NoError(t, errC)
	}

	f, err := os.Open(golden)
	require.NoError(t, err)
	t.Cleanup(func() { _ = f.Close() })

	var expected []ModulePublic

	err = json.NewDecoder(f).Decode(&expected)
	require.NoError(t, err)

	assert.Equal(t, expected, updates)
}
