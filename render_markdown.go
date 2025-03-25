package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ldez/modupwiz/internal"
	"github.com/olekukonko/tablewriter"
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
