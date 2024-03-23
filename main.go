package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	started := time.Now()

	file, err := os.Open("measurements.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	readyQueue := make(chan *worker, NUM_WORKERS)
	workers := make([]*worker, NUM_WORKERS)

	workerWg := sync.WaitGroup{}
	workerWg.Add(NUM_WORKERS)

	for i := 0; i < NUM_WORKERS; i++ {
		workers[i] = &worker{}
		go workers[i].work(readyQueue, &workerWg)
	}

	readerWg := sync.WaitGroup{}
	readerWg.Add(NUM_READERS)

	for i := 0; i < NUM_READERS; i++ {
		go reader(file, readyQueue, &readerWg)
	}

	readerWg.Wait()
	log.Printf("Readers done in %v\n", time.Since(started))

	for i := range workers {
		workers[i].workerChan <- nil
	}

	workerWg.Wait()
	log.Printf("Workers done in %v\n", time.Since(started))

	sumarize(workers)
	log.Printf("Sumarize done in %v\n", time.Since(started))

	fmt.Printf("Full process done in %v\n", time.Since(started))
}
