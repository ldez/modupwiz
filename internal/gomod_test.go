package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_extractModuleInfo(t *testing.T) {
	testCases := []struct {
		desc     string
		filename string
		expected []ModInfo
	}{
		{
			desc:     "module",
			filename: "module.json",
			expected: []ModInfo{
				{
					Path:      "github.com/ldez/modupwiz",
					Dir:       "/home/go/src/github.com/ldez/modupwiz",
					GoMod:     "/home/go/src/github.com/ldez/modupwiz/go.mod",
					GoVersion: "1.22.0",
					Main:      true,
				},
			},
		},
		{
			desc:     "workspace",
			filename: "workspace.json",
			expected: []ModInfo{
				{
					Path:      "golang.org/x/example",
					Dir:       "/home/ldez/sources/golangci-lint/modinfo/testdata/workspace/example",
					GoMod:     "/home/ldez/sources/golangci-lint/modinfo/testdata/workspace/example/go.mod",
					GoVersion: "1.15",
					Main:      true,
				},
				{
					Path:      "example.com/hello",
					Dir:       "/home/ldez/sources/golangci-lint/modinfo/testdata/workspace/hello",
					GoMod:     "/home/ldez/sources/golangci-lint/modinfo/testdata/workspace/hello/go.mod",
					GoVersion: "1.20",
					Main:      true,
				},
				{
					Path:      "example.com/world",
					Dir:       "/home/ldez/sources/golangci-lint/modinfo/testdata/workspace/world",
					GoMod:     "/home/ldez/sources/golangci-lint/modinfo/testdata/workspace/world/go.mod",
					GoVersion: "1.20",
					Main:      true,
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			file, err := os.Open(filepath.Join("fixtures", test.filename))
			require.NoError(t, err)
			t.Cleanup(func() { _ = file.Close() })

			info, err := extractModuleInfo(file)
			require.NoError(t, err)

			assert.Equal(t, test.expected, info)
		})
	}
}
