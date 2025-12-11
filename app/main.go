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



func main() {
	// TODO: Uncomment the code below to pass the first stage

	rl, err := InitReadline()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer rl.Close()

	for {

		command, err := rl.Readline()

		if err == readline.ErrInterrupt {
			continue
		}

		if err == io.EOF {
			break
		}

		// command = `"exe with \'single quotes\'" /tmp/cow/f3`
		if err = execCommand(command); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
