package main

import (
	"bytes"
	"regexp"
	"sort"
	"strings"
	"time"
)

func sortedKeys(m map[time.Time]LoadtestReport) []time.Time {
	keys := make([]time.Time, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	return keys
}

func getLabel(date time.Time) string {
	return date.Format("2006-01-02")
}

func extractKeysAndLabels(reports map[time.Time]LoadtestReport) ([]string, []string) {
	var keys []string
	var labels []string
	keyMap := map[string]bool{}
	labelMap := map[string]bool{}
	for _, key := range sortedKeys(reports) {
		for servicekey := range reports[key].Metrics {
			if !keyMap[servicekey] {
				keys = append(keys, servicekey)
				keyMap[servicekey] = true
			}
		}
		label := getLabel(reports[key].StartedAt)
		if !labelMap[label] {
			labels = append(labels, label)
			labelMap[label] = true
		}
	}
	sort.Strings(keys)
	sort.Strings(labels)
	return keys, labels
}

func formatPath(path string) string {
	match, _ := regexp.MatchString("(\\*\\.[a-zA-Z0-9]+|[*/])$", path)
	if !match {
		path = path + "/*"
	}
	if strings.HasSuffix(path, "/") {
		path = path + "*"
	}
	return path
}

func makeFirstLowerCase(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}

	bts := []byte(s)

	lc := bytes.ToLower([]byte{bts[0]})
	rest := bts[1:]

	return string(bytes.Join([][]byte{lc, rest}, nil))
}
