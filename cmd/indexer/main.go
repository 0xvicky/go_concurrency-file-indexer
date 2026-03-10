package main

import (
	"concurrent-file-indexer/internal/config"
	"concurrent-file-indexer/internal/scanner"
	"concurrent-file-indexer/internal/storage"
	"concurrent-file-indexer/internal/worker"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	println("Concurrent File Indexer")
	start := time.Now()
	dir := config.DummyFolder
	nWorkers := config.WorkerCount
	hstorage := storage.HashStorage{Hashes: make(map[string]string)}

	// //channel to pass filepath to workers
	pathChannel := make(chan string)
	//create waitGroup
	var wg sync.WaitGroup
	//spawing workers
	for i := 1; i <= nWorkers; i++ {
		wg.Add(1)
		go worker.StartWorker(i, pathChannel, &wg, &hstorage)
	}

	go scanner.Start(dir, pathChannel)

	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Time taken with %d workers:%s", nWorkers, elapsed)
	//check if map is got updated with 5000 hashes
	mapLen := len(hstorage.Hashes)
	fmt.Println("Total Files Indexed:", mapLen)
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
Filesystem✅
    ↓
Directory Scanner ✅
    ↓
File Path Channel ✅
    ↓
Worker Pool (goroutines) ✅
    ↓
Hash Computation
    ↓
Results Map (mutex protected)
    ↓
Program Summary
*/
