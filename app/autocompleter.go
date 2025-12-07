package main

import "github.com/chzyer/readline"

func InitReadline() (*readline.Instance, error) {
	return readline.NewEx(&readline.Config{
		Prompt:       prompt,
		AutoComplete: pc,
	})
}

var (
	pc = readline.NewPrefixCompleter(
		readline.PcItem("echo"),
		readline.PcItem("exit"),
	)

	prompt = "$ "
)
