package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"
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
	parsedArgs, redirectFileName := parseArgString(command)
	var f *os.File
	var err error
	if redirectFileName != "" {
		f, err = createFile(redirectFileName)
		// fmt.Println("Redirect Filename:", redirectFileName)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	w := out
	if f != nil {
		w = f
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

			if path, ok, _ := isExecutable(cmd); ok {
				fmt.Fprintf(w, "%s is %s\n", cmd, path)
				return nil
			}

			return fmt.Errorf("%s: not found", parsedArgs[1])
		}

	case "cat":
		// printSliceWithIndexAndVal(parsedArgs)
		// fmt.Println("Redirect Filename:", redirectFileName)

		for _, fname := range parsedArgs[1:] {
			fHandle, err := os.Open(fname)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cat: %s: No such file or directory\n", fname)
				continue
			}

			if f != nil {
				io.Copy(f, fHandle)
			} else {
				io.Copy(out, fHandle)
			}

			// io.Copy(w, fHandle)
			fHandle.Close()
		}

		// catCmd := exec.Command(parsedArgs[0], parsedArgs[1:]...)

		// if f != nil {
		// 	catCmd.Stdout = f
		// } else {
		// 	catCmd.Stdout = out
		// }
		// catCmd.Stderr = os.Stderr
		// err := catCmd.Run()
		// if err != nil {
		// 	// fmt.Printf("Error: %v\n", err)
		// 	return fmt.Errorf("cat: %s: No such file or directory\nError:%v\n", parsedArgs[1], err)
		// }

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
			return fmt.Errorf("cd: %s: No such file or directory", parsedArgs[1])
		}
	default:

		if _, ok, err := isExecutable(parsedArgs[0]); ok {
			// fmt.Printf("Executing:%s\nRedirect Filename: %s", parsedArgs[0], redirectFileName)
			// prepare the command to execute
			// printSliceWithIndexAndVal(parsedArgs)
			command := exec.Command(parsedArgs[0], parsedArgs[1:]...)

			if f != nil {
				command.Stdout = f
			} else {
				command.Stdout = os.Stdout
			}
			command.Stderr = os.Stderr
			// command.Stdout = out

			return command.Run()
		} else if err != nil {
			// fmt.Printf("Executing: %s\n", parsedArgs[0])
			return fmt.Errorf("%s: command not found", command)
		}
		
	}
	return nil
}

func isExecutable(cmd string) (string, bool, error) {

	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", false, err
	}
	return path, true, nil
}

func parseArgString(args string) (parsedArgs []string, redirectFileName string) {
	runes := []rune(args)
	n := len(runes)

	// var parsedArgs []string

	isSingleQuote := false
	isDoubleQuote := false
	isRedirect := false
	skipNextSpaces := false

	var rd strings.Builder
	var b strings.Builder

	for i := 0; i < n; i++ {
		r := runes[i]

		if r == '"' && !isSingleQuote {
			isDoubleQuote = !isDoubleQuote
			continue
		}

		if r == '\'' && !isDoubleQuote {
			isSingleQuote = !isSingleQuote
			continue
		}

		if r == '\\' {

			if i == n-1 {

				if isRedirect {
					rd.WriteRune('\\')
				} else {
					b.WriteRune('\\')
				}
				continue
			}

			next := runes[i+1]

			if isSingleQuote {
				b.WriteRune('\\')
				continue
			}

			if isDoubleQuote {
				// only certain escapes in double quotes
				if next == '"' || next == '\\' || next == '$' || next == '`' {
					b.WriteRune(next)
					i++
					continue
				}
				// backslash is literal otherwise
				b.WriteRune('\\')
				continue
			}

			if isRedirect {
				rd.WriteRune(next)
			} else {
				b.WriteRune(next)
			}
			i++
			continue
		}

		if !isSingleQuote && !isDoubleQuote {
			if r == '>' || (unicode.IsDigit(r) && i+1 < n && runes[i+1] == '>') {
				if b.Len() > 0 {
					parsedArgs = append(parsedArgs, b.String())
					b.Reset()
				}
				isRedirect = true
				skipNextSpaces = true
				if unicode.IsDigit(r) {
					i++
				}
				continue
			}

		}

		// consider everything after redirect as redirect file name
		if isRedirect {
			if skipNextSpaces && r == ' ' {
				continue
			}
			skipNextSpaces = false
			if r == ' ' && !isSingleQuote && !isDoubleQuote {
				isRedirect = false
				continue
			}

			rd.WriteRune(r)
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

	redirectFileName = strings.TrimSpace(rd.String())
	return parsedArgs, redirectFileName
}

func printSliceWithIndexAndVal(s []string) {
	for i, v := range s {
		log.Println(i, v)
	}
}

func createFile(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	return file, nil
}
