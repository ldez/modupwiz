package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"slices"

	"github.com/ldez/modupwiz/internal"
	"github.com/urfave/cli/v2"
	"golang.org/x/mod/modfile"
)

const (
	flgDirect   = "direct"
	flgExplicit = "explicit"
	flgIndirect = "indirect"
	flgCompare  = "compare"
	flgVersions = "versions"
	flgPipe     = "pipe"
	flgPath     = "path"
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
			Name:  flgDirect,
			Usage: "Only direct modules",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  flgExplicit,
			Usage: "Only explicit indirect modules",
		},
		&cli.BoolFlag{
			Name:  flgIndirect,
			Usage: "All indirect modules",
		},
		&cli.BoolFlag{
			Name:  flgCompare,
			Usage: "Display compare links from GitHub if possible",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  flgVersions,
			Usage: "Display versions of modules",
			// Value: true,
		},
		&cli.BoolFlag{
			Name:  flgPipe,
			Usage: "Allow to pipe the command",
		},
		&cli.StringFlag{
			Name:  flgPath,
			Usage: "File path to write the output. (Default: os.Stdout)",
		},
	}

	app.Action = func(cliCtx *cli.Context) error {
		opts := Options{
			Direct:           cliCtx.Bool(flgDirect),
			ExplicitIndirect: cliCtx.Bool(flgExplicit),
			AllIndirect:      cliCtx.Bool(flgIndirect),
			Compare:          cliCtx.Bool(flgCompare),
			Versions:         cliCtx.Bool(flgVersions),
			Pipe:             cliCtx.Bool(flgPipe),
			Path:             cliCtx.String(flgPath),
		}

		return run(context.Background(), opts)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, opts Options) error {
	goMod, err := internal.ReadGoMod(ctx)
	if err != nil {
		return fmt.Errorf("reading module file: %w", err)
	}

	updates, err := internal.ListUpdates(ctx, opts.Pipe)
	if err != nil {
		return fmt.Errorf("listing updates: %w", err)
	}

	err = render(ctx, opts, goMod, updates)
	if err != nil {
		return fmt.Errorf("rendering: %w", err)
	}

	return nil
}

func render(_ context.Context, opts Options, goMod *modfile.File, updates []internal.ModulePublic) error {
	return renderMarkdown(opts, filter(updates, opts, goMod))
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
