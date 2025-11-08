package main

import (
	"bufio"
	"fmt"
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
		code, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		os.Exit(code)
	case "echo":
		fmt.Fprintln(os.Stdout, strings.Join(args[1:], " "))
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
		// fmt.Println("Invalid command case")
		return fmt.Errorf("%s: command not found", command)
	}
	return nil
}

func isExecutable(cmd string) (string, bool) {

	path, err := exec.LookPath(cmd) // look for executable in PATH env variable
	if err != nil {
		return "", false
	}
	return path, true
}
