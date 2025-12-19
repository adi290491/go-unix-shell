package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var pnt = fmt.Fprint
var out = os.Stdout
var rl *readline.Instance

func main() {
	// TODO: Uncomment the code below to pass the first stage
	var err error

	rl, err = InitReadline()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// defer rl.Close()

	defer func() {
		if err := saveHistoryToFile(); err != nil {
			fmt.Fprintln(os.Stderr, "Warning: could not save history:", err)
		}
		rl.Close()
	}()

	for {

		command, err := rl.Readline()

		if err == readline.ErrInterrupt {
			continue
		}

		if err == io.EOF {
			break
		}
		SessionHistory = append(SessionHistory, command)

		// command = `"exe with \'single quotes\'" /tmp/cow/f3`
		if err = execCommand(command); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
