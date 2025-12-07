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
			commands: collectCommands(),
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

type ShellCompleter struct {
	commands []string
}

func (c *ShellCompleter) Do(line []rune, pos int) ([][]rune, int) {
	word, _ := getCurrentWord(line, pos)

	matches := []string{}

	for _, cmd := range c.commands {
		if strings.HasPrefix(cmd, word) {
			matches = append(matches, cmd)
		}
	}

	if len(matches) == 0 {
		fmt.Print("\x07")
		return nil, 0
	}

	out := make([][]rune, len(matches))
	for i, m := range matches {
		suffix := m[len(word):] + " "
		out[i] = []rune(suffix)
	}

	// replace := -(pos - start)
	return out, 0

}

func getCurrentWord(line []rune, pos int) (string, int) {
	i := pos - 1
	for i >= 0 && line[i] != ' ' {
		i--
	}
	return string(line[i+1 : pos]), i + 1
}

func collectCommands() []string {
	builtins := []string{"echo", "exit", "cd", "cat", "pwd", "type"}
	path := os.Getenv("PATH")
	dirs := strings.Split(path, ":")

	m := map[string]bool{}

	for _, b := range builtins {
		m[b] = true
	}

	// external commands
	for _, dir := range dirs {
		files, _ := os.ReadDir(dir)
		for _, f := range files {
			if !f.IsDir() {
				m[f.Name()] = true
			}
		}
	}

	// convert map to output slice
	var output []string

	for c := range m {
		output = append(output, c)
	}

	sort.Strings(output)

	return output
}
