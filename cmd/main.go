package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"mapgenie/analysis"
	"mapgenie/changes"
	"mapgenie/pkg/log"
	"os"
)

const description = `To define a mapper for generation write a func with the source struct in the argument 
and the destination struct in return values, then add the magic comment: 'map this pls':

type A struct {
	Count int
}

type B struct {
	Count int
}

// MyMapper map this pls
func MyMapper(A) B {
	return B{}
}

Run mapgenie on current package.

$ mapgenie ./mypkg/...

The tool will generate the body of this function. Fields with Equal names will be assigned.

If the fields have different types they might be converted, currently supported:
	* Pointers to (and from) values.
	* Values with the same underlying type.
	* Any numbers to numbers.
	* string to (and from) []byte, also any derived types.
	* string to (and from) number types (10 base).
	* time.Time to (and from):
		* string (currently only RFC3339 is supported).
		* int,int64,uint,uint64 representing seconds since the UNIX epoch (uint lie).

	* Supporting custom mapper functions is in the plans.
`

func main() {
	logLevel := log.Warn

	app := &cli.App{
		Name:        "mapgenie",
		Usage:       "Generate code for converting one struct to another",
		Description: description,
		Args:        true,
		ArgsUsage:   "./path/... [path2 path3 ...]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Set log level to debug",
				Action: func(context2 *cli.Context, debug bool) error {
					logLevel = log.Debug
					return nil
				},
			},
		},

		Action: func(c *cli.Context) error {
			c.Context = log.Ctx(c.Context, logLevel, os.Stderr)

			if c.Args().Len() == 0 {
				log.Errorf(c.Context, "Provide package paths in arguments, for example: ./...")
				cli.ShowAppHelpAndExit(c, 1)
			}

			targetFiles := analysis.FindTargetsInPackages(c.Context, c.Args().Slice()...)
			if len(targetFiles) == 0 {
				log.Warnf(c.Context, "No mappers found")
				return nil
			}

			changes.ApplyFilesChanges(c.Context, targetFiles)

			return nil
		},
	}

	ctx := context.Background()
	ctx = log.Ctx(ctx, logLevel, os.Stderr)
	err := app.RunContext(ctx, os.Args)

	if err != nil {
		log.Errorf(ctx, "%s", err.Error())
	}
}
