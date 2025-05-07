package fileutil

import (
	conf "FileFootprintLister/src/configuration"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func FormatName(name string) string {
	rules := conf.ReadRules()
	// get the last element from the path
	if hasPathSeparators(name) {
		_, name = GetFileAndPath(name)
	}

	for _, rule := range rules {
		name = replaceChar(rule, name)
	}
	return name
}

// hasPathSeparators checks if the given string contains path separators.
func hasPathSeparators(path string) bool {
	return strings.Contains(path, string(os.PathSeparator)) || strings.Contains(path, "/")
}

func replaceChar(rules []string, text string) string {
	// Define the regex pattern to match email addresses
	pattern := rules[0]
	replace := rules[1]

	// Compile the regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return text
	}

	// Replace all matches with a placeholder
	replacedText := re.ReplaceAllString(text, replace)

	return replacedText
}
