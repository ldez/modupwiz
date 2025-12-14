package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ldez/modupwiz/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_generate(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip()
	}

	modules := readJSONFile[[]internal.ModulePublic](t, filepath.Join("testdata", "modules.json"))
	sccs := readJSONFile[[][]internal.ModulePublic](t, filepath.Join("testdata", "sccs.json"))

	testCases := []struct {
		desc     string
		opts     Options
		expected string
	}{
		{
			desc: "default template (sccs)",
			opts: Options{
				Template: "#DEFAULT#",
			},
			expected: "default.golden.sh",
		},
		{
			desc: "custom template (modules)",
			opts: Options{
				Template: "./pr.md.tmpl",
			},
			expected: "pr.golden.md",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			temp, err := os.CreateTemp(t.TempDir(), test.expected)
			require.NoError(t, err)

			err = generate(test.opts, temp, modules, sccs)
			require.NoError(t, err)

			require.NoError(t, temp.Close())

			golden, err := os.ReadFile(filepath.Join("testdata", test.expected))
			require.NoError(t, err)

			result, err := os.ReadFile(temp.Name())
			require.NoError(t, err)

			assert.Equal(t, string(golden), string(result))
		})
	}
}

func readJSONFile[T any](t *testing.T, filename string) T {
	t.Helper()

	file, err := os.Open(filename)
	require.NoError(t, err)

	var a T

	err = json.NewDecoder(file).Decode(&a)
	require.NoError(t, err)

	return a
}
