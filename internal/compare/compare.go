package compare

import "sort"

func Intersection(left, right []string) []string {
	set := make(map[string]struct{}, len(right))
	for _, item := range right {
		set[item] = struct{}{}
	}

	matches := make([]string, 0)
	seen := make(map[string]struct{})
	for _, item := range left {
		if _, ok := set[item]; !ok {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		matches = append(matches, item)
	}
	sort.Strings(matches)
	return matches
}
