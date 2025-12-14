package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ldez/modupwiz/internal"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

func renderMarkdown(opts Options, modules []internal.ModulePublic) error {
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

	table := tablewriter.NewTable(writer,
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Borders: tw.Border{
				Left:      tw.On,
				Right:     tw.On,
				Top:       tw.Off,
				Bottom:    tw.Off,
				Overwrite: true,
			},
			Settings: tw.Settings{
				Separators: tw.Separators{
					ShowHeader:     tw.On,
					ShowFooter:     tw.Off,
					BetweenRows:    tw.Off,
					BetweenColumns: 0,
				},
			},
			Symbols: tw.NewSymbols(tw.StyleMarkdown),
		})))

	defer func() { _ = table.Close() }()

	titles := []string{"Module"}

	if opts.Versions {
		titles = append(titles, "Current", "Latest")
	}

	if opts.Compare {
		titles = append(titles, "Compare")
	}

	table.Header(titles)

	for _, module := range modules {
		row := []string{module.Path}

		if opts.Versions {
			row = append(row, module.Version, module.NewVersion())
		}

		if opts.Compare {
			row = append(row, getCompareLink(module))
		}

		err := table.Append(row)
		if err != nil {
			return fmt.Errorf("append row: %w", err)
		}
	}

	err := table.Render()
	if err != nil {
		return fmt.Errorf("render table: %w", err)
	}

	return nil
}
