package main

import (
	"fmt"
	"lemin/cmd/admin"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./lem-in input_file")
		os.Exit(1)
	}
	filename := os.Args[1]
	admin.Run(filename)
}
