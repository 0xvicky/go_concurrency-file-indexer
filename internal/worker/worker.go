package worker

import (
	"concurrent-file-indexer/internal/storage"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
)

func StartWorker(workerId int, workQueue <-chan string, wg *sync.WaitGroup, hs *storage.HashStorage) {

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
		hs.AddHash(finalHash, filePath)
	}
}
