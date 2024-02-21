package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(text string) []string {
	if len(text) == 0 {
		return nil
	}
	re := regexp.MustCompile(`(\s+|\s|[\^,]|[\^.][\^-])`)
	text = re.ReplaceAllString(text, " ")
	words := strings.Fields(text)
	countMap := map[string]int{}
	for _, word := range words {
		_, ok := countMap[word]
		if ok {
			countMap[word]++
		} else {
			countMap[word] = 1
		}
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
	return keys[:10]
}
