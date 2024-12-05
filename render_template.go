package main

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/ldez/modupwiz/internal"
	"golang.org/x/mod/modfile"
)

//go:embed default.sh.tmpl
var defaultTemplate string

func renderTemplate(ctx context.Context, opts Options, file *modfile.File, modules []internal.ModulePublic) error {
	if len(modules) == 0 {
		log.Println("No updates available. " + opts.scope())
		return nil
	}

	writer := os.Stdout

	if opts.Path != "" {
		var err error
		writer, err = os.Create(opts.Path)
		if err != nil {
			return fmt.Errorf("create file %s: %w", opts.Path, err)
		}

		defer func() { _ = writer.Close() }()
	}

	cmd := exec.CommandContext(ctx, "go", "mod", "graph")

	output, err := cmd.CombinedOutput()
	if err != nil {
		extra := string(output)

		var ee *exec.ExitError
		if errors.As(err, &ee) {
			extra += string(ee.Stderr)
		}

		return fmt.Errorf("command '%s': %w: %s", strings.Join(cmd.Args, " "), err, extra)
	}

	root := internal.ModulePublic{Path: file.Module.Mod.Path}

	g := internal.NewGraph()

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		a := trimAfter(parts[0])
		b := trimAfter(parts[1])

		bMod, ok := findMod(modules, b)
		if !ok {
			continue
		}

		if a == root.Path {
			g.AddEdge(root, bMod)
			continue
		}

		aMod, ok := findMod(modules, a)
		if !ok {
			continue
		}

		g.AddEdge(aMod, bMod)
	}

	sccs := g.FindSCCs()

	return generate(opts, writer, modules, removeRoot(sccs, root))
}

func removeRoot(sccs [][]internal.ModulePublic, root internal.ModulePublic) [][]internal.ModulePublic {
	var cleaned [][]internal.ModulePublic

	for _, scc := range sccs {
		var n []internal.ModulePublic

		for _, module := range scc {
			if module.Path == root.Path {
				continue
			}

			n = append(n, module)
		}

		if len(n) > 0 {
			cleaned = append(cleaned, n)
		}
	}

	return cleaned
}

func generate(opts Options, writer io.Writer, modules []internal.ModulePublic, sccs [][]internal.ModulePublic) error {
	base := template.New("default.sh.tmpl").
		Funcs(map[string]any{
			"getCompareLink": getCompareLink,
		})

	var tmpl *template.Template
	if opts.Template == "#DEFAULT#" {
		var err error
		tmpl, err = base.Parse(defaultTemplate)
		if err != nil {
			return fmt.Errorf("parse default template: %w", err)
		}
	} else {
		var err error
		tmpl, err = base.ParseFiles(opts.Template)
		if err != nil {
			return fmt.Errorf("parse template %s: %w", opts.Template, err)
		}
	}

	err := tmpl.Execute(writer, map[string]any{
		"SCCS":    sccs,
		"modules": modules,
	})
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}

func findMod(modules []internal.ModulePublic, name string) (internal.ModulePublic, bool) {
	for _, module := range modules {
		if module.Path == name {
			return module, true
		}
	}

	return internal.ModulePublic{}, false
}

func trimAfter(v string) string {
	index := strings.Index(v, "@")

	if index == -1 {
		return v
	}

	return v[:index]
}
