package output

import (
	"fmt"
	"io"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

func WriteText(w io.Writer, report Report, opts TextOptions) error {
	if opts.Quiet {
		for _, match := range report.Matches {
			if _, err := fmt.Fprintln(w, match); err != nil {
				return err
			}
		}
		return nil
	}

	if _, err := fmt.Fprintf(w, "Source: %s\n", report.Source); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Installed AUR packages checked: %d\n", report.CheckedPackageCount); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Compromised package names loaded: %d\n\n", report.CompromisedPackageCount); err != nil {
		return err
	}

	if report.Safe {
		message := "No compromised installed packages found"
		if opts.Color {
			message = colorGreen + message + colorReset
		}
		_, err := fmt.Fprintln(w, message)
		return err
	}

	heading := "WARNING: compromised packages detected:"
	if opts.Color {
		heading = colorYellow + heading + colorReset
	}
	if _, err := fmt.Fprintln(w, heading); err != nil {
		return err
	}
	for _, match := range report.Matches {
		line := "- " + match
		if opts.Color {
			line = colorRed + line + colorReset
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}
