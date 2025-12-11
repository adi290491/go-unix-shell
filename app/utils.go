package main

import (
	"fmt"
	"os"
	"strings"
)

func createFile(fileName string, redirectType RedirectType) (*os.File, error) {
	var flag int

	switch redirectType {
	case RedirectStdout, RedirectStderr:
		// '>' and '1>' and "2>"
		flag = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	case RedirectAppend, RedirectStderrAppend:
		flag = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	default:
		return nil, fmt.Errorf("invalid redirect type")
	}

	file, err := os.OpenFile(fileName, flag, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func pressRingBell() {
	fmt.Print("\x07")
}

func findLongestCommonPrefix(matches []string) string {
	if len(matches) == 0 {
		return ""
	}

	longestCommonPrefix := matches[0]

	for _, match := range matches[1:] {
		for !strings.HasPrefix(match, longestCommonPrefix) {
			if len(longestCommonPrefix) == 0 {
				return ""
			}
			longestCommonPrefix = longestCommonPrefix[:len(longestCommonPrefix)-1]
		}
	}
	return longestCommonPrefix
}
