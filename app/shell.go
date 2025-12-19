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

var builtinCmds = map[string]bool{
	"echo":    true,
	"exit":    true,
	"type":    true,
	"pwd":     true,
	"cd":      true,
	"history": true,
}

func execCommand(command string) error {
	command = strings.TrimSpace(command)

	if command == "" {
		return nil
	}

	if strings.Contains(command, "|") {

		return execPipeline(command)
	}

	return execSingleCommand(command, os.Stdin, os.Stdout, os.Stderr)
}

func execPipeline(command string) error {

	commandParts := strings.Split(command, "|")

	type Pipe struct {
		reader *os.File
		writer *os.File
	}

	type CommandInfo struct {
		cmd       string
		args      []string
		isBuiltin bool
		stdin     io.Reader
		stdout    io.Writer
		stderr    io.Writer
	}

	commands := make([]CommandInfo, 0)

	for _, part := range commandParts {
		cmd := strings.TrimSpace(part)
		args, _, _ := parseArgString(cmd)

		if len(args) == 0 {
			continue
		}

		isBuiltin := builtinCmds[args[0]]
		commands = append(commands, CommandInfo{cmd: cmd, args: args, isBuiltin: isBuiltin})
	}

	if len(commands) == 0 {
		return nil
	}

	numPipes := len(commands) - 1
	pipes := make([]Pipe, numPipes)

	// create all pipes
	for i := 0; i < numPipes; i++ {
		reader, writer, err := os.Pipe()
		if err != nil {
			return fmt.Errorf("failed to create pipe %d: %w", i, err)
		}
		pipes[i] = Pipe{reader: reader, writer: writer}
	}

	// assign pipes to each command
	for i := 0; i < len(commands); i++ {
		if i == 0 {
			commands[i].stdin = os.Stdin
		} else {
			commands[i].stdin = pipes[i-1].reader
		}

		if i == len(commands)-1 {
			commands[i].stdout = os.Stdout
		} else {
			commands[i].stdout = pipes[i].writer
		}
		commands[i].stderr = os.Stderr
	}

	errChannels := make([]chan error, len(commands))
	externalCmds := make([]*exec.Cmd, len(commands))

	// start all commands
	for i := 0; i < len(commands); i++ {
		cmd := commands[i]
		isLastCmd := i == len(commands)-1
		if cmd.isBuiltin {
			errChan := make(chan error, 1)
			errChannels[i] = errChan

			go func(c CommandInfo, ch chan error, isLast bool) {
				defer func() {
					if !isLast {
						if closer, ok := c.stdout.(io.Closer); ok {
							closer.Close()
						}
					}
				}()

				errChan <- execSingleCommand(
					c.cmd,
					c.stdin,
					c.stdout,
					c.stderr,
				)
			}(cmd, errChan, isLastCmd)

		} else {

			if _, err := exec.LookPath(cmd.args[0]); err == nil {
				execCmd := exec.Command(cmd.args[0], cmd.args[1:]...)

				if i == 0 {
					execCmd.Stdin = os.Stdin
				} else {
					execCmd.Stdin = cmd.stdin
				}

				if isLastCmd {
					execCmd.Stdout = os.Stdout
				} else {
					execCmd.Stdout = cmd.stdout
				}
				execCmd.Stderr = cmd.stderr

				externalCmds[i] = execCmd
			} else {
				fmt.Fprintf(os.Stderr, "%s: command not found\n", cmd.args[0])
			}
		}
	}

	// start all external commands
	for _, cmd := range externalCmds {
		if cmd == nil {
			continue
		}

		cmd.Start()
	}

	// close all pipe writes in parent
	for _, p := range pipes {
		p.writer.Close()
	}

	var firstErr error

	// wait for all builtin commands
	for i := 0; i < len(commands); i++ {
		if commands[i].isBuiltin {
			if err := <-errChannels[i]; err != nil {
				if firstErr == nil {
					firstErr = err
				}
			}
		}
	}

	// wait for all external commands
	for _, cmd := range externalCmds {
		if cmd == nil {
			continue
		}
		if err := cmd.Wait(); err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				continue
			}
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	for _, p := range pipes {
		p.reader.Close()
	}

	// return firstErr
	return nil
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
		// fmt.Fprintf(os.Stderr, "Execution cat: %+v\n", parsedArgs)
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
	case "history":

		if len(parsedArgs) == 1 {
			histories, err := fetchAllHistory(historyFilePath)
			if err != nil {
				return err
			}
			for _, h := range histories {
				fmt.Fprintf(w, "    %d  %s\n", h.id, h.cmd)
			}
			return nil
		}

		if len(parsedArgs) == 2 {
			lastN, err := strconv.Atoi(parsedArgs[1])
			if err != nil {
				return err
			}
			histories, err := fetchLastNHistory(lastN)
			if err != nil {
				return err
			}
			for _, h := range histories {
				fmt.Fprintf(w, "    %d  %s\n", h.id, h.cmd)
			}
			return nil
		}

		switch parsedArgs[1] {
		case "-r":

			fileName := parsedArgs[2]
			// fmt.Printf("[DEBUG] fileName: %s\n", fileName)
			file, err := os.Open(fileName)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if line != "" {
					rl.SaveHistory(line)
				}
			}

			if err := scanner.Err(); err != nil {
				return err
			}

		case "-w":
			fileName := parsedArgs[2]
			// check if file exists
			var file *os.File
			var err error

			file, err = os.Create(fileName)
			if err != nil {
				return err
			}

			defer file.Close()

			histories, err := fetchAllHistory(historyFilePath)
			if err != nil {
				return err
			}

			fileWriter := bufio.NewWriter(file)

			for _, h := range histories {
				fmt.Fprintf(fileWriter, "%s\n", h.cmd)
			}

			fileWriter.Flush()
		case "-a":
			fileName := parsedArgs[2]

			fileHistories, err := fetchAllHistory(fileName)
			if err != nil {
				return err
			}
			fileCommandSet := make(map[string]bool)
			for _, h := range fileHistories {
				fileCommandSet[h.cmd] = true
			}

			var newCommands []string

			for i := lastAppendedIdx; i < len(SessionHistory); i++ {
				h := SessionHistory[i]

				if !fileCommandSet[h] || strings.HasPrefix(h, "history -a") {
					newCommands = append(newCommands, h)
				}
			}


			// check if file exists
			var file *os.File

			file, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}

			defer file.Close()

			fileWriter := bufio.NewWriter(file)
			for _, cmd := range newCommands {
				fmt.Fprintf(fileWriter, "%s\n", cmd)
			}
			fileWriter.Flush()

			lastAppendedIdx = len(SessionHistory)
		}

	default:
		// fmt.Fprintf(os.Stderr, "Execution External: %v", parsedArgs)
		_, err := exec.LookPath(parsedArgs[0])
		var command *exec.Cmd
		if err == nil {

			command = exec.Command(parsedArgs[0], parsedArgs[1:]...)
			command.Stdin = stdin

			if f != nil && (redirectType == RedirectStdout || redirectType == RedirectAppend) {
				command.Stdout = f
			} else {
				command.Stdout = stdout
			}

			if f != nil && (redirectType == RedirectStderr || redirectType == RedirectStderrAppend) {
				command.Stderr = f
			} else {
				command.Stderr = stderr
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
