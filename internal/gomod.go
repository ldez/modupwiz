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
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/mod/modfile"
)

// ReadGoMod read the `go.mod` file.
func ReadGoMod(ctx context.Context) (*modfile.File, error) {
	modPath, err := getModulePath(ctx)
	if err != nil {
		return nil, err
	}

	return readModuleFile(modPath)
}

func getModulePath(ctx context.Context) (string, error) {
	modulePath, err := getModulePathFromEnv(ctx)
	if err != nil {
		return "", err
	}

	if modulePath != "" {
		return modulePath, nil
	}

	info, err := findModuleInfo(ctx)
	if err != nil {
		return "", fmt.Errorf("finding module info: %w", err)
	}

	return info.GoMod, nil
}

func getModulePathFromEnv(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "go", "env", "GOMOD")

	out, err := cmd.Output()
	if err != nil {
		extra := string(out)

		var ee *exec.ExitError
		if errors.As(err, &ee) {
			extra += string(ee.Stderr)
		}

		return "", fmt.Errorf("command '%s': %w: %s", strings.Join(cmd.Args, " "), err, extra)
	}

	return strings.TrimSpace(string(out)), nil
}

type ModInfo struct {
	Path      string `json:"Path"`
	Dir       string `json:"Dir"`
	GoMod     string `json:"GoMod"`
	GoVersion string `json:"GoVersion"`
	Main      bool   `json:"Main"`
}

func getModuleInfo(ctx context.Context) ([]ModInfo, error) {
	cmd := exec.CommandContext(ctx, "go", "list", "-m", "-json")

	out, err := cmd.Output()
	if err != nil {
		extra := string(out)

		var ee *exec.ExitError
		if errors.As(err, &ee) {
			extra += string(ee.Stderr)
		}

		return nil, fmt.Errorf("command '%s': %w: %s", strings.Join(cmd.Args, " "), err, extra)
	}

	info, err := extractModuleInfo(bytes.NewBuffer(out))
	if err != nil {
		return nil, fmt.Errorf("extract module info: %s", string(out))
	}

	return info, nil
}

func extractModuleInfo(in io.Reader) ([]ModInfo, error) {
	var infos []ModInfo

	for dec := json.NewDecoder(in); dec.More(); {
		var v ModInfo
		if err := dec.Decode(&v); err != nil {
			return nil, fmt.Errorf("unmarshaling error: %w", err)
		}

		if v.GoMod == "" {
			return nil, errors.New("working directory is not part of a module")
		}

		if !v.Main || v.Dir == "" {
			continue
		}

		infos = append(infos, v)
	}

	if len(infos) == 0 {
		return nil, errors.New("go.mod file not found")
	}

	// sort name length of the name: first is longer.
	sort.Slice(infos, func(i, j int) bool {
		return len(infos[i].Path) > len(infos[j].Path)
	})

	return infos, nil
}

func findModuleInfo(ctx context.Context) (ModInfo, error) {
	info, err := getModuleInfo(ctx)
	if err != nil {
		return ModInfo{}, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return ModInfo{}, fmt.Errorf("current working directory: %w", err)
	}

	var found []ModInfo
	for _, modInfo := range info {
		if strings.Contains(wd, modInfo.Dir) {
			found = append(found, modInfo)
		}
	}

	switch len(found) {
	case 0:
		return ModInfo{}, errors.New("no module information found")
	case 1:
		return found[0], nil
	default:
		return ModInfo{}, fmt.Errorf("found %d modules information, it's the root of a workspace", len(found))
	}
}

func readModuleFile(modPath string) (*modfile.File, error) {
	raw, err := os.ReadFile(filepath.Clean(modPath))
	if err != nil {
		return nil, fmt.Errorf("reading go.mod file: %w", err)
	}

	file, err := modfile.Parse("go.mod", raw, nil)
	if err != nil {
		return nil, fmt.Errorf("parse go.mod file: %w", err)
	}

	return file, nil
}
