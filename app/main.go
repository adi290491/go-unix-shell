package main

import (
	"bufio"
	"fmt"
	"io"
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
		// command = `"exe with \'single quotes\'" /tmp/cow/f3`
		if err = execCommand(command); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func execCommand(command string) error {
	command = strings.TrimSuffix(command, "\n")

	// args := strings.Split(command, " ")
	parsedArgs, redirectFileName, redirectType := parseArgString(command)

	var f *os.File
	var err error
	if redirectFileName != "" {
		f, err = createFile(redirectFileName, redirectType)
		// fmt.Println("Redirect Filename:", redirectFileName)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	w := out
	e := os.Stderr
	if f != nil {
		switch redirectType {
		case RedirectStdout:
			// fmt.Println("RedirectStdout:", redirectType)
			w = f
		case RedirectAppend:
			// fmt.Println("RedirectAppend:", redirectType)
			w = f
		case RedirectStderr:
			// fmt.Println("RedirectStderr:", redirectType)
			e = f
		}
	}

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
		// printSliceWithIndexAndVal(parsedArgs)
		outputStrings := strings.Join(parsedArgs[1:], " ")

		fmt.Fprintln(w, outputStrings)

	case "type":
		if ok := commandSets[parsedArgs[1]]; ok {
			fmt.Fprintf(w, "%s is a shell builtin\n", parsedArgs[1])
		} else {

			cmd := parsedArgs[1]

			path, err := isExecutable(cmd)

			if err == nil {
				fmt.Fprintf(w, "%s is %s\n", cmd, path)
				return nil
			} else {
				fmt.Fprintf(e, "%s: not found\n", parsedArgs[1])
			}

			// return fmt.Errorf("%s: not found", parsedArgs[1])

		}

	case "cat":
		// printSliceWithIndexAndVal(parsedArgs)
		// fmt.Println("Redirect Filename:", redirectFileName)

		for _, fname := range parsedArgs[1:] {
			fHandle, err := os.Open(fname)
			if err != nil {
				fmt.Fprintf(e, "cat: %s: No such file or directory\n", fname)
				continue
			}

			// if f != nil {
			// 	io.Copy(f, fHandle)
			// } else {
			// 	io.Copy(out, fHandle)
			// }

			io.Copy(w, fHandle)
			fHandle.Close()
		}

	case "pwd":
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%s\n", pwd)

	case "cd":

		if parsedArgs[1] == "~" {
			parsedArgs[1] = os.Getenv("HOME")
		}

		err := os.Chdir(parsedArgs[1])

		if err != nil {
			fmt.Fprintf(e, "cd: %s: No such file or directory\n", parsedArgs[1])
		}
	default:

		_, err := isExecutable(parsedArgs[0])
		var command *exec.Cmd
		if err == nil {
			// fmt.Printf("Executing: %s\n", parsedArgs[0])
			// printSliceWithIndexAndVal(parsedArgs)
			command = exec.Command(parsedArgs[0], parsedArgs[1:]...)

			if f != nil && redirectType == RedirectStdout || redirectType == RedirectAppend {
				command.Stdout = f
			} else {
				command.Stdout = os.Stdout
			}

			if f != nil && redirectType == RedirectStderr {
				command.Stderr = f
			} else {
				command.Stderr = os.Stderr
			}

			err := command.Run()
			if _, ok := err.(*exec.ExitError); ok {
				return nil
			} else {
				return err
			}
		} else {
			// fmt.Printf("Executing: %s\n", parsedArgs[0])
			fmt.Fprintf(e, "%s: command not found\n", parsedArgs[0])
		}

	}
	return nil
}

func isExecutable(cmd string) (string, error) {

	return exec.LookPath(cmd)

}

func createFile(fileName string, redirectType RedirectType) (*os.File, error) {
	var flag int

	switch redirectType {
	case RedirectStdout, RedirectStderr:
		// '>' and '1>' and "2>"
		flag = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	case RedirectAppend:
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
