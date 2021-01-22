// Package interp implements the tree-walker
// interpreter for the lox language
package interp

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	exDataErr = 65
)

// hadError records if an error was encountered earlier.
var hadError bool

// RunFile runs the lox interpreter on the
// script in the file
func RunFile(filename string) {

	script, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("unable to read ", filename)
		os.Exit(exDataErr)
	}
	run(string(script))
	if hadError {
		os.Exit(exDataErr)
	}
}

// RunPrompt runs the lox interpreter interactively
func RunPrompt() {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			fmt.Println("")
			break
		}
		run(scanner.Text())
		hadError = false
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

// error raises an error during interpretation.
func raise(line int, message string) {
	report(line, "", message)
}

// report reports an error during interpretation
func report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s",
		line, where, message)
	hadError = true
}
