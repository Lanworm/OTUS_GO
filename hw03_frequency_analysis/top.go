package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var re = regexp.MustCompile(`(\s+|\s|[\^,]|[\^.][\^-])`)

func Top10(text string) []string {
	if len(text) == 0 {
		return nil
	}
	text = re.ReplaceAllString(text, " ")
	words := strings.Fields(text)
	countMap := map[string]int{}

	for _, word := range words {
		countMap[word]++
	}
	keys := make([]string, 0, len(countMap))
	for k := range countMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if countMap[keys[i]] == countMap[keys[j]] {
			diff := []string{keys[i], keys[j]}
			sort.Strings(diff)
			return diff[0] == keys[i]
		}
		return countMap[keys[i]] > countMap[keys[j]]
	})
	if len(keys) < 10 {
		return keys
	}
	return keys[:10]
}
