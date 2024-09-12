package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_formatLinkPattern(t *testing.T) {
	module := ModulePublic{
		Version: "v0.1.0",
		Update: &ModulePublic{
			Version: "v0.2.0",
		},
	}

	testCases := []struct {
		desc     string
		name     string
		prefix   string
		expected string
	}{
		{
			desc:     "GitHub",
			name:     "github.com/example/module",
			expected: "https://github.com/example/module/compare/v0.1.0...v0.2.0",
		},
		{
			desc:     "Gitlab",
			name:     "gitlab.com/example/module",
			expected: "https://gitlab.com/example/module/-/compare/v0.1.0...v0.2.0",
		},
		{
			desc:     "Bitbucket",
			name:     "bitbucket.org/example/module",
			expected: "https://bitbucket.org/bitbucket.org/example/module/branches/compare/v0.2.0%0Dv0.1.0#diff",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			pattern := formatLinkPattern(test.name, test.prefix, module)

			assert.Equal(t, test.expected, pattern)
		})
	}
}

func Test_hasCompareLink(t *testing.T) {
	testCases := []struct {
		desc   string
		name   string
		assert assert.BoolAssertionFunc
	}{
		{
			desc:   "GitHub",
			name:   "github.com/example/module",
			assert: assert.True,
		},
		{
			desc:   "Gitlab",
			name:   "gitlab.com/example/module",
			assert: assert.True,
		},
		{
			desc:   "Bitbucket",
			name:   "bitbucket.org/example/module",
			assert: assert.True,
		},
		{
			desc:   "Other",
			name:   "example.org/example/module",
			assert: assert.False,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			test.assert(t, hasCompareLink(test.name))
		})
	}
}

func Test_getCompareLink(t *testing.T) {
	testCases := []struct {
		desc     string
		module   ModulePublic
		expected string
	}{
		{
			desc: "GitHub",
			module: ModulePublic{
				Path:    "github.com/example/module",
				Version: "v0.1.0",
				Update: &ModulePublic{
					Path:    "github.com/example/module",
					Version: "v0.2.0",
				},
			},
			expected: "https://github.com/example/module/compare/v0.1.0...v0.2.0",
		},
		{
			desc: "Gitlab",
			module: ModulePublic{
				Path:    "gitlab.com/example/module",
				Version: "v0.1.0",
				Update: &ModulePublic{
					Path:    "gitlab.com/example/module",
					Version: "v0.2.0",
				},
			},
			expected: "https://gitlab.com/example/module/-/compare/v0.1.0...v0.2.0",
		},
		{
			desc: "Bitbucket",
			module: ModulePublic{
				Path:    "bitbucket.org/example/module",
				Version: "v0.1.0",
				Update: &ModulePublic{
					Path:    "bitbucket.org/example/module",
					Version: "v0.2.0",
				},
			},
			expected: "https://bitbucket.org/bitbucket.org/example/module/branches/compare/v0.2.0%0Dv0.1.0#diff",
		},
		{
			desc: "major version",
			module: ModulePublic{
				Path:    "github.com/example/module/v2",
				Version: "v0.1.0",
				Update: &ModulePublic{
					Path:    "github.com/example/module/v2",
					Version: "v0.2.0",
				},
			},
			expected: "https://github.com/example/module/compare/v0.1.0...v0.2.0",
		},
		{
			desc: "submodule",
			module: ModulePublic{
				Path:    "github.com/example/module/test",
				Version: "v0.1.0",
				Update: &ModulePublic{
					Path:    "github.com/example/module/test",
					Version: "v0.2.0",
				},
			},
			expected: "https://github.com/example/module/compare/test/v0.1.0...test/v0.2.0",
		},
		{
			desc: "submodule and major version",
			module: ModulePublic{
				Path:    "github.com/example/module/test/v2",
				Version: "v0.1.0",
				Update: &ModulePublic{
					Path:    "github.com/example/module/test/v2",
					Version: "v0.2.0",
				},
			},
			expected: "https://github.com/example/module/compare/test/v0.1.0...test/v0.2.0",
		},
		{
			desc: "golang.org/x",
			module: ModulePublic{
				Path:    "golang.org/x/crypto",
				Version: "v0.1.0",
				Update: &ModulePublic{
					Path:    "golang.org/x/crypto",
					Version: "v0.2.0",
				},
			},
			expected: "https://github.com/golang/crypto/compare/v0.1.0...v0.2.0",
		},
		{
			desc: "vanity and submodule",
			module: ModulePublic{
				Path:    "cloud.google.com/go/compute",
				Version: "v0.1.0",
				Update: &ModulePublic{
					Path:    "cloud.google.com/go/compute",
					Version: "v0.2.0",
				},
			},
			expected: "https://github.com/GoogleCloudPlatform/gcloud-golang/compare/compute/v0.1.0...compute/v0.2.0",
		},
		{
			desc: "vanity and long submodule",
			module: ModulePublic{
				Path:    "go.elastic.co/apm/module/apmhttp",
				Version: "v0.1.0",
				Update: &ModulePublic{
					Path:    "go.elastic.co/apm/module/apmhttp",
					Version: "v0.2.0",
				},
			},
			expected: "https://github.com/elastic/apm-agent-go/compare/module/apmhttp/v0.1.0...module/apmhttp/v0.2.0",
		},
		{
			desc: "vanity and without submodule but long URL",
			module: ModulePublic{
				Path:    "go4.org/unsafe/assume-no-moving-gc",
				Version: "v0.1.0",
				Update: &ModulePublic{
					Path:    "go4.org/unsafe/assume-no-moving-gc",
					Version: "v0.2.0",
				},
			},
			expected: "https://github.com/go4org/unsafe-assume-no-moving-gc/compare/v0.1.0...v0.2.0",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			link := getCompareLink(test.module)

			assert.Equal(t, test.expected, link)
		})
	}
}
