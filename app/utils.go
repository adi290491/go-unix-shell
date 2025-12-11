package main

import (
	"fmt"
	"os"
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
