package worker

import (
	"concurrent-file-indexer/internal/storage"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"sync"
)

func StartWorker(workerId int, workQueue <-chan string, wg *sync.WaitGroup, hs *storage.HashStorage) {

	defer wg.Done()
	for filePath := range workQueue {
		// println("Worker id", workerId, "processing", filePath)
		file, err := os.Open(filePath)
		if err != nil {
			println("Fault in filepath", err)
			continue
		}

		chunk := make([]byte, 1024) //buffer to store the chunk
		chunkHasher := sha256.New()
		for {
			chunkSize, err := file.Read(chunk)
			if chunkSize > 0 {
				chunkHasher.Write(chunk[0:chunkSize])
				// println(chunkSize, hex.EncodeToString(hash[:]))
			}

			//check if EOF
			if err != nil {
				if err == io.EOF {
					// println("End of the current file")
					break
				}
			}

		}
		file.Close()

		fileHash := hex.EncodeToString(chunkHasher.Sum(nil))
		// fmt.Println(fileHash)
		hs.AddHash(fileHash, filePath)
	}
}
