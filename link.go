package main

import (
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"

	"github.com/ldez/grignotin/metago"
)

const (
	repoNameLength = 2
	repoNameParts  = 3
)

func getCompareLink(module ModulePublic) string {
	name := module.Path

	if module.NewVersion() == "" {
		return ""
	}

	var prefix string
	if !hasCompareLink(name) {
		switch {
		case strings.HasPrefix(name, "golang.org/x/"):
			// custom override to use GitHub instead of https://cs.opensource.google
			name = strings.ReplaceAll(name, "golang.org/x/", "github.com/golang/")

		default:
			meta, err := metago.Get(name)
			if err != nil {
				log.Println(err)
				return module.Path
			}

			name = metago.EffectivePkgSource(meta)
			prefix = strings.TrimPrefix(strings.TrimPrefix(meta.Pkg, meta.GoImport[0]), "/")

			if !hasCompareLink(name) {
				return ""
			}
		}
	}

	if prefix == "" && strings.Count(name, "/") > repoNameLength {
		name, prefix = splitName(name)
	}

	return formatLinkPattern(name, prefix, module)
}

func hasCompareLink(name string) bool {
	return strings.HasPrefix(name, "github.com/") || strings.HasPrefix(name, "gitlab.com/") || strings.HasPrefix(name, "bitbucket.org/")
}

func splitName(name string) (string, string) {
	n := strings.Split(name, "/")

	repo := path.Join(n[:repoNameParts]...)

	if ok, _ := regexp.MatchString(`v\d+`, n[len(n)-1]); ok {
		if len(n) <= repoNameParts {
			return repo, ""
		}

		return repo, path.Join(n[repoNameParts : len(n)-1]...)
	}

	return repo, path.Join(n[repoNameParts:]...)
}

func formatLinkPattern(name, prefix string, module ModulePublic) string {
	var pattern string
	switch {
	case strings.HasPrefix(name, "github.com/"):
		pattern = "https://%s/compare/%s...%s"
	case strings.HasPrefix(name, "gitlab.com/"):
		pattern = "https://%s/-/compare/%s...%s"
	case strings.HasPrefix(name, "bitbucket.org/"):
		pattern = "https://bitbucket.org/%s/branches/compare/%[3]s%%0D%[2]s#diff"
	}

	return fmt.Sprintf(pattern,
		name,
		path.Join(prefix, extractVersion(module.Version)),
		path.Join(prefix, extractVersion(module.NewVersion())),
	)
}

func extractVersion(v string) string {
	expVersion := regexp.MustCompile(`v.+[.-]\d{14}-[a-z0-9]{12}`)

	switch {
	case strings.HasSuffix(v, "+incompatible"):
		return strings.TrimSuffix(v, "+incompatible")

	case strings.HasPrefix(v, "v0.0.0-") || expVersion.MatchString(v):
		return v[strings.LastIndex(v, "-")+1:]

	default:
		return v
	}
}
