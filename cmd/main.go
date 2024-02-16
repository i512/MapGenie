package main

import (
	"context"
	"flag"
	"mapgenie/analysis"
	"mapgenie/changes"
	"mapgenie/pkg/log"
	"os"
)

func main() {
	ctx := context.Background()
	ctx = log.Ctx(ctx, log.Warn, os.Stderr)

	flag.Parse()
	patterns := flag.Args()
	if len(patterns) == 0 {
		log.Errorf(ctx, "Provide package patterns in arguments, for example: ./...")
	}

	targetFiles := analysis.FindTargetsInPackages(ctx, patterns...)
	changes.ApplyFilesChanges(ctx, targetFiles)
}
