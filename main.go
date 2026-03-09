package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

// producer function.
func scanner(dir string, workQueue chan<- string) {
	println("Scanner")
	/*
		walk through a directory✅
		find files ✅
		send file paths into the channel✅
	*/
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		// println(path)
		workQueue <- path
	}
	close(workQueue)

}

func worker(workerId int, workQueue <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for filePath := range workQueue {
		println("Worker id", workerId, "processing", filePath)
	}
}

func main() {
	println("Concurrent File Indexer")
	dir := "Z:/Code/Golang/concurrent-file-indexer/dummy"

	//map to store result
	// m := make(map[string]string)
	// //channel to pass filepath to workers
	pathChannel := make(chan string)
	go scanner(dir, pathChannel)
	//create waitGroup
	var wg sync.WaitGroup

	//spawing workers
	nWorkers := 3
	for i := 1; i <= nWorkers; i++ {
		wg.Add(1)
		go worker(i, pathChannel, &wg)
	}
	wg.Wait()

}

/*
Error 1
An unbuffered channel requires both sides to be ready at the same time:

sender  ↔  receiver

Think of it like handing someone a package. You cannot drop it unless someone is there to receive it.

Right now the situation is:

scanner → trying to send
main → not receiving yet

Because main() is still waiting for scanner() to return.

So the send operation blocks.

The scanner waits forever for someone to receive.

But main() will only start receiving after scanner() finishes.

And scanner() cannot finish because it’s stuck sending.

You now have a perfect circular wait.

That is the deadlock.

=================================================

*/

/*
Program starts
      ↓
Create shared structures (channels, result map, mutex, waitgroup)
      ↓
Start worker goroutines
      ↓
Scanner walks through directory
      ↓
Scanner sends file paths into channel
      ↓
Workers continuously read file paths
      ↓
Workers open file and compute hash
      ↓
Workers safely store result
      ↓
Scanner finishes → channel closed
      ↓
Workers detect closed channel → exit
      ↓
WaitGroup waits for all workers
      ↓
Program prints results
      ↓
Program exits
*/

/*
Filesystem
    ↓
Directory Scanner
    ↓
File Path Channel
    ↓
Worker Pool (goroutines)
    ↓
Hash Computation
    ↓
Results Map (mutex protected)
    ↓
Program Summary
*/
