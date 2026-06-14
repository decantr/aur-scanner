package local

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

func ReadAURPackages(ctx context.Context, pacmanPath string) ([]string, error) {
	cmd := exec.CommandContext(ctx, pacmanPath, "-Qqm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return nil, fmt.Errorf("pacman not found in PATH")
		}
		return nil, fmt.Errorf("%w: %s", err, strings.TrimSpace(string(output)))
	}
	return parsePackageLines(output), nil
}

func parsePackageLines(data []byte) []string {
	seen := make(map[string]struct{})
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		pkg := strings.TrimSpace(scanner.Text())
		if pkg == "" {
			continue
		}
		seen[pkg] = struct{}{}
	}

	packages := make([]string, 0, len(seen))
	for pkg := range seen {
		packages = append(packages, pkg)
	}
	sort.Strings(packages)
	return packages
}
