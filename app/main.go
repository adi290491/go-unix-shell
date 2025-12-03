package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var pnt = fmt.Fprint
var out = os.Stdout

var commandSets = map[string]bool{
	"echo": true,
	"exit": true,
	"type": true,
	"pwd":  true,
	"cd":   true,
}

func main() {
	// TODO: Uncomment the code below to pass the first stage
	for {
		fmt.Fprint(os.Stdout, "$ ")
		reader := bufio.NewReader(os.Stdin)

		command, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err = execCommand(command); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func execCommand(command string) error {
	command = strings.TrimSuffix(command, "\n")

	args := strings.Split(command, " ")
	parsedArgs := parseArgString(strings.Join(args, " "))

	switch parsedArgs[0] {
	case "exit":
		if len(parsedArgs) == 1 {
			os.Exit(0)
		}
		code, err := strconv.Atoi(parsedArgs[1])
		if err != nil {
			return err
		}
		os.Exit(code)
	case "echo":

		outputStrings := strings.Join(parsedArgs[1:], " ")
		fmt.Fprintln(os.Stdout, outputStrings)
	case "type":
		if ok := commandSets[parsedArgs[1]]; ok {
			fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", parsedArgs[1])
		} else {

			cmd := parsedArgs[1]

			if path, ok := isExecutable(cmd); ok {
				fmt.Fprintf(out, "%s is %s\n", cmd, path)
				return nil
			}

			return fmt.Errorf("%s: not found", parsedArgs[1])
		}
	case "cat":

		copyCmd := exec.Command(parsedArgs[0], parsedArgs[1:]...)
		copyCmd.Stderr = os.Stderr
		copyCmd.Stdout = out

		return copyCmd.Run()
	case "pwd":
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%s\n", pwd)
	case "cd":

		if parsedArgs[1] == "~" {
			parsedArgs[1] = os.Getenv("HOME")
		}

		err := os.Chdir(parsedArgs[1])

		if err != nil {
			return fmt.Errorf("cd: %s: No such file or directory", parsedArgs[1])
		}
	default:


		if _, ok := isExecutable(parsedArgs[0]); ok {

			// prepare the command to execute
			command := exec.Command(parsedArgs[0], parsedArgs[1:]...)

			command.Stderr = os.Stderr
			command.Stdout = out

			return command.Run()
		}
		return fmt.Errorf("%s: command not found", command)
	}
	return nil
}

func isExecutable(cmd string) (string, bool) {

	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", false
	}
	return path, true
}

func parseArgString(args string) []string {
	runes := []rune(args)
	n := len(runes)

	var parsedArgs []string

	isSingleQuote := false
	isDoubleQuote := false
	var b strings.Builder

	for i := 0; i < n; i++ {
		r := runes[i]

		if r == '\\' {
			// inside single quotes: backslash is literal
			if isSingleQuote {
				b.WriteRune('\\')
				continue
			}

			// if last char, backslash is literal
			if i == n-1 {
				b.WriteRune('\\')
				continue
			}

			next := runes[i+1]
			if isDoubleQuote {
				// inside double quotes, backslash only escapes ", \, $, `
				if next == '"' || next == '\\' || next == '$' || next == '`' {
					b.WriteRune(next)
					i++
					continue
				} else {
					b.WriteRune('\\')
					continue
				}
			}

			// Outside of any quotes: backslash escapes the next char
			b.WriteRune(next)
			i++
			continue
		}

		if r == '"' && !isSingleQuote {
			isDoubleQuote = !isDoubleQuote
			continue
		}

		if r == '\'' && !isDoubleQuote {
			isSingleQuote = !isSingleQuote
			continue
		}

		if r == ' ' && !isSingleQuote && !isDoubleQuote {
			if b.Len() > 0 {
				parsedArgs = append(parsedArgs, b.String())
				b.Reset()
			}
			continue
		}

		b.WriteRune(r)
	}

	if b.Len() > 0 {
		parsedArgs = append(parsedArgs, b.String())
	}
	return parsedArgs
}

func printSliceWithIndexAndVal(s []string) {
	for i, v := range s {
		log.Println(i, v)
	}
}

