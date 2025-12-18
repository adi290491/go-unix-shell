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

/*
	// for i := 0; i < len(commands); i++ {
	// 	// isLast := (i == len(commands)-1)

	// 	if commands[i].isBuiltin {
	// 		err := <-errChannels[i]
	// 		// fmt.Fprintf(os.Stderr, "Builtin %d finished: %v\n", i, err)
	// 		if err != nil && firstErr == nil {
	// 			// fmt.Fprintf(os.Stderr, "[DEBUG] [BUILTIN] error is %v\n", err)
	// 			firstErr = err
	// 		}
	// 	} else {
	// 		if externalCmds[i] != nil {
	// 			// fmt.Fprintf(os.Stderr, "Waiting for external %s\n", commands[i].args[0])
	// 			err := externalCmds[i].Wait()

	// 			if err != nil {
	// 				// fmt.Fprintf(os.Stderr, "[DEBUG] externalCmds[%v] error is %v\n", externalCmds[i], err)
	// 				if _, ok := err.(*exec.ExitError); ok {
	// 					// fmt.Fprintf(os.Stderr, "error is exit error\n")
	// 					continue
	// 				}

	// 				if firstErr == nil {
	// 					firstErr = err
	// 				}
	// 			}

	// 		} else {
	// 			fmt.Fprintf(os.Stderr, "[DEBUG] externalCmds[%d] is NIL!\n", i)
	// 		}
	// 	}
	// }

*/