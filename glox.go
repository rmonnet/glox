package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	exUsage   = 64
	exDataErr = 65
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
	run(string(script))
}

// runPrompt runs the lox interpreter interactively
func runPrompt() {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			fmt.Println("")
			break
		}
		run(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error while reading ", err)
		os.Exit(exDataErr)
	}

}

// run runs the lox interpreter on the provided
// program.
func run(script string) {

	fmt.Println(script)
}
