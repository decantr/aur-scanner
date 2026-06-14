package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"aur-scanner/internal/compare"
	"aur-scanner/internal/local"
	"aur-scanner/internal/output"
	"aur-scanner/internal/source"
)

const DefaultURL = "https://md.archlinux.org/s/SxbqukK6IA/download"

type scanOptions struct {
	URL   string
	File  string
	Quiet bool
}

func RunScan(ctx context.Context, stdout io.Writer, stderr io.Writer, args []string) int {
	opts, err := parseScanOptions(stderr, args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		fmt.Fprintf(stderr, "ERROR: %v\n", err)
		return 2
	}

	installedPackages, err := local.ReadAURPackages(ctx, "pacman")
	if err != nil {
		fmt.Fprintf(stderr, "ERROR: failed to read installed AUR packages: %v\n", err)
		return 2
	}

	effectiveSource := opts.URL
	var sourceData []byte

	if opts.File != "" {
		effectiveSource = opts.File
		sourceData, err = source.ReadFile(opts.File)
	} else {
		effectiveSource = source.NormalizeURL(opts.URL)
		sourceData, err = source.FetchURL(ctx, effectiveSource)
	}

	if err != nil {
		fmt.Fprintf(stderr, "ERROR: failed to load compromised package source: %v\n", err)
		return 2
	}

	compromisedPackages, err := source.ParseCompromisedPackages(sourceData)
	if err != nil {
		fmt.Fprintf(stderr, "ERROR: failed to parse compromised package source: %v\n", err)
		return 2
	}

	matches := compare.Intersection(installedPackages, compromisedPackages)
	report := output.Report{
		Source:                  effectiveSource,
		CheckedPackageCount:     len(installedPackages),
		CompromisedPackageCount: len(compromisedPackages),
		Matches:                 matches,
		Safe:                    len(matches) == 0,
	}

	if err := output.WriteText(stdout, report, output.TextOptions{Quiet: opts.Quiet}); err != nil {
		fmt.Fprintf(stderr, "ERROR: failed to write output: %v\n", err)
		return 2
	}

	if len(matches) > 0 {
		return 1
	}
	return 0
}

func parseScanOptions(stderr io.Writer, args []string) (scanOptions, error) {
	opts := scanOptions{}
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&opts.URL, "url", DefaultURL, "URL to fetch compromised package data from")
	fs.StringVar(&opts.File, "file", "", "read compromised package data from a local file")
	fs.BoolVar(&opts.Quiet, "quiet", false, "print only matching package names")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "Usage:")
		fmt.Fprintln(stderr, "  aur-scanner scan [flags]")
		fmt.Fprintln(stderr, "  aur-scanner [flags]")
		fmt.Fprintln(stderr, "")
		fmt.Fprintln(stderr, "Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return opts, err
		}
		return opts, err
	}
	if opts.File != "" && opts.URL != DefaultURL {
		return opts, fmt.Errorf("--file and --url cannot be used together")
	}
	if fs.NArg() != 0 {
		return opts, fmt.Errorf("unexpected arguments: %v", fs.Args())
	}
	return opts, nil
}
