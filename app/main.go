package main

import (
	"bufio"
	"fmt"
	"os"
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
			directories := strings.Split(os.Getenv("PATH"), ":")

			for _, dir := range directories {

				path := dir + "/" + cmd

				fi, err := os.Stat(path)

				if err != nil {
					if os.IsNotExist(err) {
						continue
					}
				}

				if !fi.IsDir() {

					if fi.Mode()&0111 != 0 {
						fmt.Fprintf(out, "%s is %s\n", cmd, path)
						return nil
					}

				}
			}
			fmt.Fprintf(os.Stdout, "%s: not found\n", args[1])
		}
	default:
		// fmt.Println("Invalid command case")
		return fmt.Errorf("%s: command not found", command)
	}
	return nil
}
