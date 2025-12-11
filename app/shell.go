package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var builtinCmds = map[string]bool{
	"echo": true,
	"exit": true,
	"type": true,
	"pwd":  true,
	"cd":   true,
}

func execCommand(command string) error {
	command = strings.TrimSpace(command)

	if command == "" {
		return nil
	}

	if strings.Contains(command, "|") {
		return execPipeline(command)
	}

	return execSingleCommand(command, nil, os.Stdout, os.Stderr)
}

func execPipeline(command string) error {
	parts := strings.Split(command, "|")
	if len(parts) != 2 {
		return fmt.Errorf("only two-command pipelines are supported")
	}

	leftCmd := strings.TrimSpace(parts[0])
	rightCmd := strings.TrimSpace(parts[1])

	leftArgs := strings.Fields(leftCmd)
	rightArgs := strings.Fields(rightCmd)

	if len(leftArgs) == 0 || len(rightArgs) == 0 {
		return fmt.Errorf("empty command in pipeline")
	}

	reader, writer, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("failed to create pipe: %w", err)
	}

	// check if command is a builtin
	leftIsBuiltin := builtinCmds[leftArgs[0]]
	rightIsBuiltin := builtinCmds[rightArgs[0]]

	switch {
	case leftIsBuiltin && rightIsBuiltin:

		errChan := make(chan error, 1)
		go func() {
			err := execSingleCommand(leftCmd, os.Stdin, writer, os.Stderr)
			writer.Close()
			errChan <- err
		}()

		err := execSingleCommand(rightCmd, reader, os.Stdout, os.Stderr)
		reader.Close()

		if leftErr := <-errChan; leftErr != nil {
			return leftErr
		}
		return err

	case leftIsBuiltin && !rightIsBuiltin:

		if err := execSingleCommand(leftCmd, os.Stdin, writer, os.Stderr); err != nil {
			writer.Close()
			reader.Close()
			return err
		}

		writer.Close()

		rightExec := exec.Command(rightArgs[0], rightArgs[1:]...)
		rightExec.Stdin = reader
		rightExec.Stdout = os.Stdout
		rightExec.Stderr = os.Stderr

		err := rightExec.Run()
		reader.Close()
		return err

	case !leftIsBuiltin && rightIsBuiltin:

		leftExec := exec.Command(leftArgs[0], leftArgs[1:]...)
		leftExec.Stdin = os.Stdin
		leftExec.Stdout = writer
		leftExec.Stderr = os.Stderr

		if err := leftExec.Start(); err != nil {
			writer.Close()
			reader.Close()
			return fmt.Errorf("failed to start %s: %w", leftArgs[0], err)
		}
		writer.Close()

		if err := execSingleCommand(rightCmd, reader, os.Stdout, os.Stderr); err != nil {
			reader.Close()
			leftExec.Wait()
			return err
		}

		reader.Close()

		return leftExec.Wait()

	default:
		// setup left cmd
		leftExec := exec.Command(leftArgs[0], leftArgs[1:]...)
		leftExec.Stdin = os.Stdin
		leftExec.Stdout = writer
		leftExec.Stderr = os.Stderr

		// setup right cmd
		rightExec := exec.Command(rightArgs[0], rightArgs[1:]...)
		rightExec.Stdin = reader
		rightExec.Stdout = os.Stdout
		rightExec.Stderr = os.Stderr

		// start left command
		if err := leftExec.Start(); err != nil {
			writer.Close()
			reader.Close()
			return fmt.Errorf("failed to start %s: %w", leftArgs[0], err)
		}

		// start right command
		if err := rightExec.Start(); err != nil {
			writer.Close()
			reader.Close()
			return fmt.Errorf("failed to start %s: %w", rightArgs[0], err)
		}

		// Close write end in parent
		writer.Close()

		// wait for left command to finish
		leftExec.Wait()

		// close read end
		reader.Close()

		//wait for right command to finish
		return rightExec.Wait()
	}

}

func execSingleCommand(command string, stdin io.Reader, stdout, stderr io.Writer) error {
	command = strings.TrimSuffix(command, "\n")
	parsedArgs, redirectFileName, redirectType := parseArgString(command)

	var f *os.File
	var err error
	if redirectFileName != "" {
		f, err = createFile(redirectFileName, redirectType)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	w := stdout
	e := stderr
	if f != nil {
		switch redirectType {
		case RedirectStdout, RedirectAppend:
			w = f
		case RedirectStderr, RedirectStderrAppend:
			e = f
		}
	}

	if len(parsedArgs) == 0 {
		return nil
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
		outputStrings := strings.Join(parsedArgs[1:], " ")

		fmt.Fprintln(w, outputStrings)

	case "type":
		if len(parsedArgs) < 2 {
			fmt.Fprintln(e, "type: missing operand")
			return nil
		}

		if ok := builtinCmds[parsedArgs[1]]; ok {
			fmt.Fprintf(w, "%s is a shell builtin\n", parsedArgs[1])
		} else {

			cmd := parsedArgs[1]

			path, err := exec.LookPath(cmd)

			if err == nil {
				fmt.Fprintf(w, "%s is %s\n", cmd, path)
				return nil
			} else {
				fmt.Fprintf(e, "%s: not found\n", parsedArgs[1])
			}

		}

	case "cat":

		for _, fname := range parsedArgs[1:] {
			fHandle, err := os.Open(fname)
			if err != nil {
				fmt.Fprintf(e, "cat: %s: No such file or directory\n", fname)
				continue
			}

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

		if len(parsedArgs) == 1 {
			os.Chdir(os.Getenv("HOME"))
			return nil
		}

		if parsedArgs[1] == "~" {
			parsedArgs[1] = os.Getenv("HOME")
		}

		err := os.Chdir(parsedArgs[1])

		if err != nil {
			fmt.Fprintf(e, "cd: %s: No such file or directory\n", parsedArgs[1])
		}
	default:

		_, err := exec.LookPath(parsedArgs[0])
		var command *exec.Cmd
		if err == nil {

			command = exec.Command(parsedArgs[0], parsedArgs[1:]...)

			if f != nil && (redirectType == RedirectStdout || redirectType == RedirectAppend) {
				command.Stdout = f
			} else {
				command.Stdout = os.Stdout
			}

			if f != nil && (redirectType == RedirectStderr || redirectType == RedirectStderrAppend) {
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
