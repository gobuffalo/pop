package main

import (
	"fmt"
	"os"

	"github.com/markbates/grift/cmd"
)

func main() {
	err := cmd.Run("grift", os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
