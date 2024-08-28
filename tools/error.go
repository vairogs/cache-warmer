package tools

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintError(err error) {
	if err == nil {
		return
	}

	fmt.Println(fmt.Sprintf("%s /!\\ %s /!\\", color.New(color.FgHiRed).Sprint("/!\\"), err))
}
