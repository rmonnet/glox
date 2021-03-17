package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rmonnet/glox/interp"
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

	parseOnly := flag.Bool("parseOnly", false, "parse and dump the AST")
	flag.Parse()
	args := flag.Args()

	if len(args) > 1 {
		fmt.Println("Usage glox [-parseOnly] [script]")
		os.Exit(exUsage)
	} else if len(args) == 1 {
		runFile(args[0], *parseOnly)
	} else {
		runPrompt(*parseOnly)
	}
}

// runFile runs the lox interpreter on the
// script in the file
func runFile(filename string, parseOnly bool) {

	script, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("unable to read ", filename)
		os.Exit(exDataErr)
	}
	interp := interp.New(os.Stdout, os.Stderr)
	interp.Run(string(script), parseOnly)
	if interp.HadCompileError() {
		os.Exit(exDataErr)
	}
	if interp.HadRuntimeError() {
		os.Exit(exSwErr)
	}
}

// runPrompt runs the lox interpreter interactively
func runPrompt(parseOnly bool) {

	scanner := bufio.NewScanner(os.Stdin)
	interp := interp.New(os.Stdout, os.Stderr)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			fmt.Println("")
			break
		}
		interp.Run(scanner.Text(), parseOnly)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error while reading ", err)
		os.Exit(exDataErr)
	}

}
