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

	// channel to pass file paths to workers
	// buffered channel with capacity 100 to decouple scanner and workers
	pathChannel := make(chan string, 100)
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
PROJECT DESCRIPTION
1. Flowchart of project
2. Description about the project working and the components built
3. Improvements that being done in the project
*/
