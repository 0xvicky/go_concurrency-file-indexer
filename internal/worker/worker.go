package worker

import (
	"concurrent-file-indexer/internal/storage"
	"fmt"
	"io"
	"os"
	"sync"
)

func StartWorker(workerId int, workQueue <-chan string, wg *sync.WaitGroup, hs *storage.HashStorage) {

	defer wg.Done()
	for filePath := range workQueue {
		// println("Worker id", workerId, "processing", filePath)
		//read the file content using filepath
		file, err := os.Open(filePath)
		if err != nil {
			println("Fault in filepath", err)
			continue
		}
		defer file.Close()
		fmt.Println(file)

		nChunks := 0

		for {
			chunk := make([]byte, 1024)
			chunkSize, err := file.Read(chunk)
			if chunkSize > 0 {
				println(chunkSize)
				nChunks++
			}

			if err != nil {
				if err == io.EOF {
					println("End of the current file")
					break
				}
			}

		}

		//calculate the hash using sha256
		// byteHash := sha256.New()
		// byteHash.Write(fileContent)
		// finalHash := hex.EncodeToString(byteHash.Sum(nil))

		// fmt.Println(finalHash)
		// hs.AddHash(finalHash, filePath)
	}
}
