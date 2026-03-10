package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type hashStorage struct {
	mu  sync.Mutex
	mpp map[string]string
}

func (h *hashStorage) addHash(resHash string, filePath string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.mpp[filePath] = resHash
}

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
func dispatcher(dir string, workQueue chan<- string) {
	// println("Scanner")
	defer close(workQueue)
	recursiveScanner(dir, workQueue)
}

func worker(workerId int, workQueue <-chan string, wg *sync.WaitGroup, hs *hashStorage) {

	defer wg.Done()
	for filePath := range workQueue {
		// println("Worker id", workerId, "processing", filePath)
		//read the file content using filepath
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// println(string(fileContent))
		//calculate the hash using sha256
		byteHash := sha256.New()
		byteHash.Write(fileContent)
		finalHash := hex.EncodeToString(byteHash.Sum(nil))
		// fmt.Println(finalHash)
		hs.addHash(finalHash, filePath)
	}
}

func main() {
	println("Concurrent File Indexer")
	start := time.Now()
	dir := "Z:/Code/Golang/concurrent-file-indexer/dummy"
	hstorage := hashStorage{mpp: make(map[string]string)}

	// //channel to pass filepath to workers
	pathChannel := make(chan string)

	//create waitGroup
	var wg sync.WaitGroup
	//spawing workers
	nWorkers := 8
	for i := 1; i <= nWorkers; i++ {
		wg.Add(1)
		go worker(i, pathChannel, &wg, &hstorage)
	}

	go dispatcher(dir, pathChannel)

	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Time taken with %d workers:%s", nWorkers, elapsed)
	//check if map is got updated with 5000 hashes
	mapLen := len(hstorage.mpp)
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
