package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

func main() {
	// TODO: Uncomment the code below to pass the first stage
	fmt.Fprint(os.Stdout, "$ ")
	reader := bufio.NewReader(os.Stdin)

	command, err := reader.ReadString('\n')

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Fprintf(os.Stdout, "%s: command not found", command[:len(command)-1])
}
