// Package interp implements the tree-walker
// interpreter for the lox language
package interp

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"gitlab.com/rcmonnet/glox/lang"
)

const (
	exDataErr = 65
)

// RunFile runs the lox interpreter on the
// script in the file
func RunFile(filename string) {

	script, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("unable to read ", filename)
		os.Exit(exDataErr)
	}
	run(string(script))
	if lang.HadError {
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
		if lang.HadError {
			os.Exit(exDataErr)
		}
		lang.HadError = false
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error while reading ", err)
		os.Exit(exDataErr)
	}

}

// run runs the lox interpreter on the provided
// program.
func run(script string) {

	scanner := lang.NewScanner(script)
	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}
}
