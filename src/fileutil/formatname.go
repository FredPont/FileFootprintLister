/*
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Written by Frederic PONT.
(c) Frederic Pont 2024
*/

package fileutil

import (
	conf "FileFootprintLister/src/configuration"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

// rulesOnce ensures the regex rules are loaded from disk exactly once,
// regardless of how many directories are processed.
var (
	rulesOnce  sync.Once
	cachedRules [][]string
)

// loadRules loads rules.csv once and caches the result.
func loadRules() [][]string {
	rulesOnce.Do(func() {
		cachedRules = conf.ReadRules()
	})
	return cachedRules
}

// FormatName applies the regex substitution rules to name and returns the
// transformed string. Path separators are stripped before applying rules.
func FormatName(name string) string {
	rules := loadRules() // cached — no disk I/O after first call

	if hasPathSeparators(name) {
		_, name = GetFileAndPath(name)
	}

	for _, rule := range rules {
		name = replaceChar(rule, name)
	}
	return name
}

// hasPathSeparators reports whether path contains OS or Unix path separators.
func hasPathSeparators(path string) bool {
	return strings.Contains(path, string(os.PathSeparator)) ||
		strings.Contains(path, "/")
}

// replaceChar applies one regex substitution rule to text.
func replaceChar(rule []string, text string) string {
	if len(rule) < 2 {
		return text
	}
	re, err := regexp.Compile(rule[0])
	if err != nil {
		fmt.Println("Regex compile error:", err)
		return text
	}
	return re.ReplaceAllString(text, rule[1])
}
