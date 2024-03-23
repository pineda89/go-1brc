package main

import (
	"github.com/edsrzf/mmap-go"
	"os"
	"sync"
)

type BufferData struct {
	n int
}

func reader(file *os.File, readyQueue chan *worker, wg *sync.WaitGroup) {
	mmap, err := mmap.Map(file, mmap.RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer mmap.Unmap()

	off := int64(0)
	for {
		worker := <-readyQueue
		if off >= int64(len(mmap)) {
			wg.Done()
			break
		}
		copy(worker.buffer, mmap[off:min(off+int64(FILE_BUFFER_SIZE), int64(len(mmap)-1))])

		n := FILE_BUFFER_SIZE
		for j := len(worker.buffer) - 1; j >= 0; j-- {
			if worker.buffer[j] == '\n' {
				off += int64(j) + 1
				n = j
				break
			}
		}

		worker.workerChan <- &BufferData{
			n: n,
		}
	}
}
