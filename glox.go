package main

import (
	"fmt"
	"os"

	"gitlab.com/rcmonnet/glox/interp"
)

const (
	exUsage = 64
)

// main runs the glox interpreter command line
// it will:
//   - interpret the script passed as argument
//   - run the lox shell if no argument is passed
//   - error if more than one argument is passed
func main() {

	if len(os.Args) > 2 {
		fmt.Println("Usage glox [script]")
		os.Exit(exUsage)
	} else if len(os.Args) == 2 {
		interp.RunFile(os.Args[1])
	} else {
		interp.RunPrompt()
	}
}
