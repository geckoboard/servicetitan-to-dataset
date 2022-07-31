package main

import (
	"os"

	"servicetitan-to-dataset/cmd"
)

func main() {
	root := cmd.Setup()

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
