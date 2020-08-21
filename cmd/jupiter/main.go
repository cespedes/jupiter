package main

import (
	"fmt"
	"os"

	"github.com/cespedes/jupiter"
)

func main() {
	fmt.Println("Starting Jupiter")
	j, err := jupiter.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i := 1; i < len(os.Args); i++ {
		score, err := j.Write(0, []byte(os.Args[i]))
		if err != nil {
			fmt.Printf("Error writing arg %d (%q): %s\n", i, os.Args[i], err)
		}
		fmt.Printf("Arg %d (%q): score=%s\n", i, os.Args[i], score)
	}
}
