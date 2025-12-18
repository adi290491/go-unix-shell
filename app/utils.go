package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
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

func dump(command string) {
	fmt.Fprintf(os.Stderr, "\n=== DEBUG PIPELINE ===\n")
	fmt.Fprintf(os.Stderr, "Command: %s\n\n", command)

	// Test each stage manually
	commandParts := strings.Split(command, "|")
	var lastOutput []byte

	for i, part := range commandParts {
		cmdStr := strings.TrimSpace(part)
		args, _, _ := parseArgString(cmdStr)

		if len(args) == 0 {
			continue
		}

		fmt.Fprintf(os.Stderr, "--- Stage %d: %s ---\n", i+1, cmdStr)

		cmd := exec.Command(args[0], args[1:]...)

		if i > 0 {
			cmd.Stdin = bytes.NewReader(lastOutput)
		}

		output, err := cmd.Output()

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			if exitErr, ok := err.(*exec.ExitError); ok {
				fmt.Fprintf(os.Stderr, "Exit code: %d\n", exitErr.ExitCode())
			}
		}

		lines := bytes.Count(output, []byte("\n"))
		if len(output) > 0 && output[len(output)-1] != '\n' {
			lines++ // Count last line if no trailing newline
		}

		fmt.Fprintf(os.Stderr, "Output: %d bytes, %d lines\n", len(output), lines)

		if len(output) > 0 {
			scanner := bufio.NewScanner(bytes.NewReader(output))
			lineNum := 1
			for scanner.Scan() {
				fmt.Fprintf(os.Stderr, "  %d: %q\n", lineNum, scanner.Text())
				lineNum++
			}
		} else {
			fmt.Fprintf(os.Stderr, "  (empty)\n")
		}

		fmt.Fprintf(os.Stderr, "\n")
		lastOutput = output
	}

	fmt.Fprintf(os.Stderr, "=== END DEBUG ===\n\n")
}
