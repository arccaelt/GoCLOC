package main

import (
	"os"
	"parser"
	"utils"
	"runtime"
)

const ERROR_OS_EXIT_CODE int = 0

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		println("The program needs a directory or file to be specified")
		os.Exit(ERROR_OS_EXIT_CODE)
	}
	inputPath := args[0]

	if !utils.IsDirectory(inputPath) {
		println("Please provide a path to the directory containing the code.")
		os.Exit(ERROR_OS_EXIT_CODE)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	parser.DirectoryAnalyzer(args[0])
	parser.DisplayResultsSummary()
}
