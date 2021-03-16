package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"gitlab.com/rcmonnet/glox/interp"
)

const (
	exUsage   = 64
	exDataErr = 65
	exSwErr   = 70
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
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

// runFile runs the lox interpreter on the
// script in the file
func runFile(filename string) {

	script, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("unable to read ", filename)
		os.Exit(exDataErr)
	}
	interp := interp.New(os.Stdout, os.Stderr)
	interp.Run(string(script))
	if interp.HadCompileError() {
		os.Exit(exDataErr)
	}
	if interp.HadRuntimeError() {
		os.Exit(exSwErr)
	}
}

// runPrompt runs the lox interpreter interactively
func runPrompt() {

	scanner := bufio.NewScanner(os.Stdin)
	interp := interp.New(os.Stdout, os.Stderr)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			fmt.Println("")
			break
		}
		interp.Run(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error while reading ", err)
		os.Exit(exDataErr)
	}

}
