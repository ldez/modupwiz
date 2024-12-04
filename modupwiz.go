package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"slices"

	"github.com/ldez/modupwiz/internal"
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
		info, err := internal.FindModuleInfo(ctx)
		if err != nil {
			return fmt.Errorf("finding module info: %w", err)
		}

		file, err = internal.ReadModuleFile(info)
		if err != nil {
			return fmt.Errorf("reading module file: %w", err)
		}
	}

	updates, err := internal.ListUpdates(ctx, opts.Pipe)
	if err != nil {
		return fmt.Errorf("listing updates: %w", err)
	}

	err = render(opts, filter(updates, opts, file))
	if err != nil {
		return fmt.Errorf("rendering: %w", err)
	}
	return nil
}

func filter(updates []internal.ModulePublic, opts Options, file *modfile.File) []internal.ModulePublic {
	var modules []internal.ModulePublic

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

func render(opts Options, modules []internal.ModulePublic) error {
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
