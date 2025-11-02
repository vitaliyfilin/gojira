package main

import (
	"fmt"
	"os"

	"gojira/internal/load"
)

func main() {
	if err := load.Run(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
