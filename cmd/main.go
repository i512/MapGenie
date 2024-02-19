package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"mapgenie/analysis"
	"mapgenie/changes"
	"mapgenie/pkg/log"
	"os"
	// "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "mapgenie",
		Usage:     "Generate code for converting one struct to another",
		Args:      true,
		ArgsUsage: "./path/... [path2 ...]",
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				log.Errorf(c.Context, "Provide package paths in arguments, for example: ./...")
				cli.ShowAppHelpAndExit(c, 1)
			}

			targetFiles := analysis.FindTargetsInPackages(c.Context, c.Args().Slice()...)
			changes.ApplyFilesChanges(c.Context, targetFiles)

			return nil
		},
	}

	ctx := context.Background()
	ctx = log.Ctx(ctx, log.Warn, os.Stderr)

	err := app.RunContext(ctx, os.Args)

	if err != nil {
		log.Errorf(ctx, "%s", err.Error())
	}
}
