package lang

import (
	"fmt"
	"os"
)

// TODO: this should really use golang errors

// HadError records if an error was encountered earlier.
var HadError bool

// Raise raises an error during interpretation.
func Raise(line int, message string) {
	report(line, "", message)
}

// report reports an error during interpretation
func report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s",
		line, where, message)
	HadError = true
}
