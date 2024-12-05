package fileutil

import (
	conf "FileFootprintLister/src/configuration"
	"fmt"
	"regexp"
)

func FormatName(name string) string {
	rules := conf.ReadRules()

	for _, rule := range rules {
		name = replaceChar(rule, name)
	}
	return name
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
