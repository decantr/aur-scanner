package main

import (
	"context"
	"os"

	"aur-scanner/cmd"
)

func main() {
	ctx := context.Background()
	args := os.Args[1:]

	if len(args) > 0 {
		switch args[0] {
		case "scan":
			os.Exit(cmd.RunScan(ctx, os.Stdout, os.Stderr, args[1:]))
		case "help":
			os.Exit(cmd.RunScan(ctx, os.Stdout, os.Stderr, []string{"--help"}))
		}
	}

	os.Exit(cmd.RunScan(ctx, os.Stdout, os.Stderr, args))
}
