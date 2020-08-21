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
	scores := []jupiter.Score{}
	fmt.Printf("Writing %d blocks...\n", len(os.Args)-1)
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("Writing %q...\n", os.Args[i])
		score, err := j.Write(0, []byte(os.Args[i]))
		if err != nil {
			fmt.Printf("Error writing arg %d (%q): %s\n", i, os.Args[i], err)
		}
		scores = append(scores, score)
		fmt.Printf("Arg %d (%q): score=%s\n", i, os.Args[i], score)
	}
	fmt.Printf("Reading %d blocks...\n", len(os.Args)-1)
	for i := 0; i < len(os.Args)-1; i++ {
		fmt.Printf("Reading %s...\n", scores[i])
		t, b, err := j.Read(scores[i])
		if err != nil {
			fmt.Printf("Error reading: %s\n", err)
		}
		fmt.Printf("score %s: type=%d len=%d data=%q\n", scores[i], t, len(b), b)
	}
}
