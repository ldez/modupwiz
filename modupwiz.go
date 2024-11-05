package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	"golang.org/x/mod/modfile"
)

var version = "dev"

// Options runner options.
type Options struct {
	Direct           bool
	ExplicitIndirect bool
	AllIndirect      bool
	Compare          bool
	Versions         bool
	Pipe             bool
	Path             string
}

func (o Options) scope() string {
	return fmt.Sprintf("(direct: %v, explicit: %v, all indirect: %v)", o.Direct, o.ExplicitIndirect, o.AllIndirect)
}

func main() {
	app := cli.NewApp()
	app.Name = "modupwiz"
	app.HelpName = "modupwiz"
	app.Usage = "Modules Update Wizard"
	app.EnableBashCompletion = true

	app.Version = version
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("modupwiz version %s %s/%s\n", c.App.Version, runtime.GOOS, runtime.GOARCH)
	}

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "direct",
			Usage: "Only direct modules",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "explicit",
			Usage: "Only explicit indirect modules",
		},
		&cli.BoolFlag{
			Name:  "indirect",
			Usage: "All indirect modules",
		},
		&cli.BoolFlag{
			Name:  "compare",
			Usage: "Display compare links from GitHub if possible",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "versions",
			Usage: "Display versions of modules",
			// Value: true,
		},
		&cli.BoolFlag{
			Name:  "pipe",
			Usage: "Allow to pipe the command",
		},
		&cli.StringFlag{
			Name:  "path",
			Usage: "File path to write the output. (Default: os.Stdout)",
		},
	}

	app.Action = func(cliCtx *cli.Context) error {
		opts := Options{
			Direct:           cliCtx.Bool("direct"),
			ExplicitIndirect: cliCtx.Bool("explicit"),
			AllIndirect:      cliCtx.Bool("indirect"),
			Compare:          cliCtx.Bool("compare"),
			Versions:         cliCtx.Bool("versions"),
			Pipe:             cliCtx.Bool("pipe"),
			Path:             cliCtx.String("path"),
		}

		return run(context.Background(), opts)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, opts Options) error {
	var file *modfile.File
	if opts.ExplicitIndirect && !opts.AllIndirect {
		info, err := findModuleInfo(ctx)
		if err != nil {
			return err
		}

		file, err = readModuleFile(info)
		if err != nil {
			return err
		}
	}

	updates, err := listUpdates(ctx, opts.Pipe)
	if err != nil {
		return err
	}

	err = renderScript(opts, filter(updates, opts, file))
	if err != nil {
		return err
	}
	return nil
}

func filter(updates []ModulePublic, opts Options, file *modfile.File) []ModulePublic {
	var modules []ModulePublic

	for _, update := range updates {
		// direct deps
		if opts.Direct && !update.Indirect {
			if update.Replace != nil {
				// Skip replaced dependencies.
				continue
			}

			modules = append(modules, update)
			continue
		}

		if update.Indirect {
			if update.Replace != nil {
				// Skip replaced dependencies.
				continue
			}

			// indirect deps
			if opts.AllIndirect {
				modules = append(modules, update)
				continue
			}

			if opts.ExplicitIndirect {
				if file == nil {
					log.Println("no data from go.mod, skip explicit indirect")
					continue
				}

				// explicit indirect deps
				foundRequire := slices.ContainsFunc(file.Require, func(require *modfile.Require) bool {
					return require.Mod.Path == update.Path
				})

				if foundRequire {
					modules = append(modules, update)
				}

				continue
			}
		}
	}

	return modules
}

func renderScript(opts Options, modules []ModulePublic) error {
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

	const (
		nameX      = "golang.org/x/"
		nameXOauth = "golang.org/x/oauth2"
		nameOther  = "other"
	)

	_, _ = fmt.Fprintln(writer, "#!/bin/sh -e")
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "git fmul")
	_, _ = fmt.Fprintln(writer, "git reset --hard upstream/master")
	_, _ = fmt.Fprintln(writer)

	groups := make(map[string][]ModulePublic)

	for _, module := range modules {
		switch {
		case strings.HasPrefix(module.Path, nameXOauth):
			groups[nameXOauth] = append(groups[nameXOauth], module)
		case strings.HasPrefix(module.Path, nameX):
			groups[nameX] = append(groups[nameX], module)
		default:
			groups[nameOther] = append(groups[nameOther], module)
		}
	}

	if len(groups[nameX]) > 0 {
		_, _ = fmt.Fprintln(writer, "echo \"chore: update golang.org/x\"")
		_, _ = fmt.Fprint(writer, "go get")
		for _, module := range groups[nameX] {
			_, _ = fmt.Fprintf(writer, " %s@%s", module.Path, module.NewVersion())
		}
		_, _ = fmt.Fprintln(writer)
		_, _ = fmt.Fprintln(writer, "go mod tidy")
		_, _ = fmt.Fprintln(writer, "git add .; git commit -m \"chore: update golang.org/x\"")
		_, _ = fmt.Fprintln(writer)
	}

	if len(groups[nameXOauth]) > 0 {
		_, _ = fmt.Fprintln(writer, "echo \"chore: update golang.org/x/oauth2\"")
		_, _ = fmt.Fprint(writer, "go get")
		for _, module := range groups[nameXOauth] {
			_, _ = fmt.Fprintf(writer, " %s@%s", module.Path, module.NewVersion())
		}
		_, _ = fmt.Fprintln(writer)
		_, _ = fmt.Fprintln(writer, "go mod tidy")
		_, _ = fmt.Fprintln(writer, "git add .; git commit -m \"chore: update golang.org/x/oauth2\"")
		_, _ = fmt.Fprintln(writer)
	}

	for _, module := range groups[nameOther] {
		_, _ = fmt.Fprintf(writer, "# %s\n", getCompareLink(module))
		_, _ = fmt.Fprintf(writer, "echo \"chore: update %s\"\n", module.Path)
		_, _ = fmt.Fprintf(writer, "go get %s@%s\n", module.Path, module.NewVersion())
		_, _ = fmt.Fprintln(writer, "go mod tidy")
		_, _ = fmt.Fprintf(writer, "git add .; git commit -m \"chore: update %s\"\n", module.Path)
		_, _ = fmt.Fprintln(writer)
	}

	return nil
}

func render(opts Options, modules []ModulePublic) error {
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

	table := tablewriter.NewWriter(writer)

	titles := []string{"Module"}
	if opts.Versions {
		titles = append(titles, "Current", "Latest")
	}
	if opts.Compare {
		titles = append(titles, "Compare")
	}

	table.SetHeader(titles)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, module := range modules {
		row := []string{module.Path}

		if opts.Versions {
			row = append(row, module.Version, module.NewVersion())
		}

		if opts.Compare {
			row = append(row, getCompareLink(module))
		}

		table.Append(row)
	}

	table.Render()

	return nil
}
