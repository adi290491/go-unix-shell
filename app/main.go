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

	switch args[0] {
	case "exit":
		if len(args) == 1 {
			os.Exit(0)
		}
		code, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		os.Exit(code)
	case "echo":
		// parsedArgs := parseArgString("test     hello")
		parsedArgs := parseArgString(strings.Join(args[1:], " "))
		outputStrings := strings.Join(parsedArgs, " ")
		fmt.Fprintln(os.Stdout, outputStrings)
	case "type":
		if ok := commandSets[args[1]]; ok {
			fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", args[1])
		} else {

			cmd := args[1]

			if path, ok := isExecutable(cmd); ok {
				fmt.Fprintf(out, "%s is %s\n", cmd, path)
				return nil
			}

			return fmt.Errorf("%s: not found", args[1])
		}
	case "cat":

		// parsedArgs := parseArgString("'/tmp/ant/f   84' '/tmp/ant/f   68' '/tmp/ant/f   25'")
		parsedArgs := parseArgString(strings.Join(args[1:], " "))

		copyCmd := exec.Command(args[0], parsedArgs...)
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

		if args[1] == "~" {
			args[1] = os.Getenv("HOME")
		}

		err := os.Chdir(args[1])

		if err != nil {
			return fmt.Errorf("cd: %s: No such file or directory", args[1])
		}
	default:

		// cmd, arguments := args[0], args[1:]

		// fmt.Printf("Command: %+v\tArgs: %+v\n", cmd, arguments)
		if _, ok := isExecutable(args[0]); ok {
			// fmt.Fprintf(out, "%s is %s\n", cmd, path)

			// prepare the command to execute
			command := exec.Command(args[0], args[1:]...)
			// fmt.Printf("Command to execute: %+v", command)
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

	parsedArgs := []string{}
	// isBackslash := false
	isSingleQuote := false
	isDoubleQuote := false
	var b strings.Builder

	for i, r := range args {
		if r == '\\' {
			// inside single quotes: backslash is literal
			if isSingleQuote {
				b.WriteRune('\\')
				continue
			}

			// if last char, backslash is literal
			if i == len(args)-1 {
				b.WriteRune('\\')
				continue
			}

			next := rune(args[i+1])
			if isDoubleQuote {
				// inside double quotes, backslash only escapes ", \, $, `
				if next == '"' || next == '\\' || next == '$' || next == '`' {
					continue
				} else {
					b.WriteRune('\\')
					continue
				}
			}

			// Outside of any quotes: backslash escapes the next char
			continue
		}
		if r == '"' {
			if i > 0 && args[i-1] == '\\' {
				b.WriteRune(r)
			} else {
				isDoubleQuote = !isDoubleQuote
			}
			continue
		}
		if r == '\'' && !isDoubleQuote {
			if i > 0 && args[i-1] == '\\' {
				b.WriteRune(r)
			} else {
				isSingleQuote = !isSingleQuote
			}
			continue
		}
		if r == ' ' && !isSingleQuote && !isDoubleQuote && (i > 0 && args[i-1] != '\\') {
			if b.Len() > 0 {
				parsedArgs = append(parsedArgs, b.String())
				b.Reset()
			}
			continue
		}

		b.WriteRune(r)
	}

	parsedArgs = append(parsedArgs, b.String())
	return parsedArgs
}

func printSliceWithIndexAndVal(s []string) {
	for i, v := range s {
		log.Println(i, v)
	}
}

/*
parsedArgs = []string{}
if quote then isQuote = !isQuote
if isQuote
  accum(c)
  continue
else
  accum(c)
  continue


*/
