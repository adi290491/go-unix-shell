package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/chzyer/readline"
)

func InitReadline() (*readline.Instance, error) {
	return readline.NewEx(&readline.Config{
		Prompt: prompt,
		AutoComplete: &ShellCompleter{
			Commands: collectCommands(),
		},
	})
}

var (
	pc = readline.NewPrefixCompleter(
		readline.PcItem("echo"),
		readline.PcItem("exit"),
	)

	prompt = "$ "
)

const (
	executable = iota
	builtin
)

// func tabInterceptor(r rune) (rune, bool) {
// 	if r != '\t' {
// 		wasLastTab = false
// 		return r, true
// 	}

// 	if !wasLastTab {
// 		wasLastTab = true
// 		pressRingBell()
// 		return 0, false
// 	}

// 	wasLastTab = false
// 	handleDoubleTab()
// 	return 0, false
// }

// func handleDoubleTab() {

// }

type ShellCompleter struct {
	Commands
	lastLine      string
	tabPressCount int
	lastPos       int
}

type Commands struct {
	Builtins    []string
	Executables []string
}

func (c *ShellCompleter) Do(line []rune, pos int) ([][]rune, int) {

	currLine := string(line[:pos])

	if c.lastLine == currLine && pos == c.lastPos {
		c.tabPressCount++
	} else {
		c.tabPressCount = 1
		c.lastLine = currLine
		c.lastPos = pos
	}

	word, _ := getCurrentWord(line, pos)

	matches := c.getAllMatches(word)

	if len(matches) == 0 {
		fmt.Print("\x07")
		c.tabPressCount = 0
		return nil, 0
	}

	if len(matches) == 1 {
		suffix := matches[0][len(word):] + " "
		c.tabPressCount = 0
		return [][]rune{[]rune(suffix)}, 0
	}

	if c.tabPressCount == 1 {
		fmt.Print("\x07")
		return nil, 0
	} else if c.tabPressCount >= 2 {
		sort.Strings(matches)
		// out := make([][]rune, len(matches))
		// for i, m := range matches {
		// 	suffix := m + " "
		// 	out[i] = []rune(suffix)
		// }
		fmt.Fprintf(os.Stderr, "\n%s\n", strings.Join(matches, "  "))
		c.tabPressCount = 0
		// return out, 0
		return [][]rune{[]rune("")}, -1
	}

	return nil, 0

}

func getCurrentWord(line []rune, pos int) (string, int) {
	i := pos - 1
	for i >= 0 && line[i] != ' ' {
		i--
	}
	return string(line[i+1 : pos]), i + 1
}

func collectCommands() Commands {
	builtins := []string{"echo", "exit", "cd", "cat", "pwd", "type"}

	builtinSet := map[string]bool{}

	for _, b := range builtins {
		builtinSet[b] = true
	}

	path := os.Getenv("PATH")
	dirs := strings.Split(path, ":")

	execSet := map[string]bool{}

	// external commands
	for _, dir := range dirs {
		files, _ := os.ReadDir(dir)
		for _, f := range files {
			if !f.IsDir() {
				execSet[f.Name()] = true
			}
		}
	}

	exec := []string{}
	for e := range execSet {
		if !builtinSet[e] {
			exec = append(exec, e)
		}
	}

	sort.Strings(builtins)
	sort.Strings(exec)

	return Commands{
		Builtins:    builtins,
		Executables: exec,
	}
}

func (c *ShellCompleter) getAllMatches(word string) []string {
	matches := make([]string, 0)

	builtins := c.getBuiltinMatches(word)
	matches = append(matches, builtins...)

	executables := c.getExecutableMatches(word)
	matches = append(matches, executables...)

	return matches
}

func (c *ShellCompleter) getBuiltinMatches(word string) []string {
	matches := make([]string, 0)

	for _, cmd := range c.Builtins {
		if strings.HasPrefix(cmd, word) {
			matches = append(matches, cmd)

		}
	}

	return matches
}

func (c *ShellCompleter) getExecutableMatches(word string) []string {
	matches := make([]string, 0)

	for _, e := range c.Executables {
		if strings.HasPrefix(e, word) {
			matches = append(matches, e)
		}
	}
	return matches
}

func pressRingBell() {
	fmt.Print("\x07")
}
