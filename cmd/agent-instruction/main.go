package main

import (
	"fmt"
	"os"

	"github.com/validkeys/agent-instruction/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
