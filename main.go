package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No input was provided")
		return
	}

	take := os.Args[1:]
	input := take[0]
	//calls banner functions
	if input == "" {
		fmt.Println()
		return
	}
}
