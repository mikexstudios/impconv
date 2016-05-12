package main

import "fmt"
import "os"
import "path/filepath"

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		fmt.Println("Usage:", filepath.Base(os.Args[0]), "[chi-data.txt]")
		os.Exit(1)
	}
	in := os.Args[1]

	fmt.Println(in)
}
