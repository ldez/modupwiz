package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ModulePublic Inspired by https://github.com/golang/go/blob/643d816c8b4348850a8a2a622d73256beea104cd/src/cmd/go/internal/modinfo/info.go
type ModulePublic struct {
	Path       string        `json:",omitempty"` // module path
	Version    string        `json:",omitempty"` // module version
	Query      string        `json:",omitempty"` // version query corresponding to this version
	Versions   []string      `json:",omitempty"` // available module versions
	Replace    *ModulePublic `json:",omitempty"` // replaced by this module
	Time       *time.Time    `json:",omitempty"` // time version was created
	Update     *ModulePublic `json:",omitempty"` // available update (with -u)
	Main       bool          `json:",omitempty"` // is this the main module?
	Indirect   bool          `json:",omitempty"` // module is only indirectly needed by main module
	Dir        string        `json:",omitempty"` // directory holding local copy of files, if any
	GoMod      string        `json:",omitempty"` // path to go.mod file describing module, if any
	GoVersion  string        `json:",omitempty"` // go version used in module
	Retracted  []string      `json:",omitempty"` // retraction information, if any (with -retracted or -u)
	Deprecated string        `json:",omitempty"` // deprecation message, if any (with -u)
	Error      *ModuleError  `json:",omitempty"` // error loading module

	Reuse bool `json:",omitempty"` // reuse of old module info is safe
}

func (m ModulePublic) NewVersion() string {
	if m.Update == nil {
		return ""
	}

	return m.Update.Version
}

type ModuleError struct {
	Err string // error text
}

func ListUpdates(ctx context.Context, pipe bool) ([]ModulePublic, error) {
	if pipe {
		return extractUpdates(os.Stdin)
	}

	cmd := exec.CommandContext(ctx, "go", "list", "-u", "-m", "-retracted", "-json", "all")

	out, err := cmd.Output()
	if err != nil {
		extra := string(out)

		var ee *exec.ExitError
		if errors.As(err, &ee) {
			extra += string(ee.Stderr)
		}

		return nil, fmt.Errorf("command '%s': %w: %s", strings.Join(cmd.Args, " "), err, extra)
	}

	updates, err := extractUpdates(bytes.NewBuffer(out))
	if err != nil {
		return nil, fmt.Errorf("extract updates: %w: %s", err, string(out))
	}

	return updates, nil
}

func extractUpdates(in io.Reader) ([]ModulePublic, error) {
	var modules []ModulePublic

	for dec := json.NewDecoder(in); dec.More(); {
		var m ModulePublic
		if err := dec.Decode(&m); err != nil {
			return nil, fmt.Errorf("unmarshaling error: %w", err)
		}

		if m.Main {
			continue
		}

		if m.Update == nil {
			continue
		}

		modules = append(modules, m)
	}

	return modules, nil
}
