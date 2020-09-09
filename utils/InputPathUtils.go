package utils

import "os"

func IsDirectory(inputPath string) bool {
	file := openInputPath(inputPath)
	stat, err := file.Stat()

	if err != nil {
		println("Can't get information about the specified path")
		os.Exit(0)
	}

	return stat.IsDir()
}

func openInputPath(inputPath string) *os.File {
	file, err := os.Open(inputPath)

	if err != nil {
		println("Invalid path")
		os.Exit(0)
	}

	return file
}
