package scanner

import (
	"fmt"
	"os"
	"path/filepath"
)

func recursiveScanner(dir string, workQueue chan<- string) {
	/*
		walk through a directory✅
		find files ✅
		if sub dir, scan them as well
		send file paths into the channel✅

	*/
	// println("hello")
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			recursiveScanner(path, workQueue)
		} else {
			// println(path)
			workQueue <- path
		}
	}
}

// producer function.
func Start(dir string, workQueue chan<- string) {
	// println("Scanner")
	defer close(workQueue)
	recursiveScanner(dir, workQueue)
}
