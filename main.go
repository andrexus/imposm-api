package main

//go:generate enumer -type=APIResponseStatus -json enums

import (
	"fmt"
	"os"

	"github.com/andrexus/imposm-api/cmd"
)

func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run command: %v\n", err)
		os.Exit(1)
	}
}
