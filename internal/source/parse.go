package source

import (
	"fmt"
	"html"
	"regexp"
	"sort"
	"strings"
)

var (
	packageTokenRE    = regexp.MustCompile(`^[a-z0-9][a-z0-9@._+-]*$`)
	markdownBodyRE    = regexp.MustCompile(`(?is)<div[^>]*class="[^"]*markdown-body[^"]*"[^>]*>(.*?)</div>`)
	metaDescriptionRE = regexp.MustCompile(`(?is)<meta[^>]*name="description"[^>]*content="([^"]*)"[^>]*>`)
	stripTagsRE       = regexp.MustCompile(`(?is)<[^>]+>`)
	fencedBlockRE     = regexp.MustCompile("(?s)```(?:[^\\n]*)\\n?(.*?)```")
)

func ParseCompromisedPackages(data []byte) ([]string, error) {
	body := string(data)
	segments := candidateSegments(body)

	for _, segment := range segments {
		if packages := parseFencedBlocks(segment); len(packages) > 0 {
			return packages, nil
		}
	}
	for _, segment := range segments {
		if packages := parsePackageText(segment); len(packages) > 0 {
			return packages, nil
		}
	}
	return nil, fmt.Errorf("no package names found in source")
}

func candidateSegments(body string) []string {
	segments := make([]string, 0, 3)
	if markdownBody := extractMarkdownBody(body); markdownBody != "" {
		segments = append(segments, markdownBody)
	}
	if metaDescription := extractMetaDescription(body); metaDescription != "" {
		segments = append(segments, metaDescription)
	}
	segments = append(segments, body)
	return segments
}

func extractMarkdownBody(body string) string {
	matches := markdownBodyRE.FindStringSubmatch(body)
	if len(matches) < 2 {
		return ""
	}
	cleaned := html.UnescapeString(matches[1])
	cleaned = stripTagsRE.ReplaceAllString(cleaned, " ")
	return cleaned
}

func extractMetaDescription(body string) string {
	matches := metaDescriptionRE.FindStringSubmatch(body)
	if len(matches) < 2 {
		return ""
	}
	return html.UnescapeString(matches[1])
}

func parseFencedBlocks(input string) []string {
	matches := fencedBlockRE.FindAllStringSubmatch(input, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]struct{})
	for _, match := range matches {
		for _, pkg := range parsePackageText(match[1]) {
			seen[pkg] = struct{}{}
		}
	}
	return sortUnique(seen)
}

func parsePackageText(input string) []string {
	seen := make(map[string]struct{})
	for _, token := range strings.Fields(input) {
		token = strings.TrimSpace(strings.Trim(token, ",;[](){}<>\"'`"))
		if !packageTokenRE.MatchString(token) {
			continue
		}
		seen[token] = struct{}{}
	}
	return sortUnique(seen)
}

func sortUnique(seen map[string]struct{}) []string {
	packages := make([]string, 0, len(seen))
	for pkg := range seen {
		packages = append(packages, pkg)
	}
	sort.Strings(packages)
	return packages
}
