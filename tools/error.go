package tools

import (
	"fmt"

	"github.com/fatih/color"
)

// PrintError prints an error message to the console. If the given error is nil, it does nothing.
// Otherwise, it formats the error message with a red warning symbol and prints it in the console.
func PrintError(err error) {
	if err == nil {
		return
	}

	fmt.Println(fmt.Sprintf("%s /!\\ %s /!\\", color.New(color.FgHiRed).Sprint("/!\\"), err))
}
