package main

import (
	"encoding/binary"
	"github.com/dolthub/swiss"
	"sync"
)

const (
	line_jump_mask = 0x0A0A0A0A0A0A0A0A
	semicolon_mask = 0x3B3B3B3B3B3B3B3B
)

type worker struct {
	swdata     *swiss.Map[uint64, *stationData]
	buffer     []byte
	workerChan chan *BufferData
}

func (w *worker) work(readyQueue chan *worker, wg *sync.WaitGroup) {
	w.buffer = make([]byte, FILE_BUFFER_SIZE)
	w.swdata = swiss.NewMap[uint64, *stationData](uint32(1024))
	w.workerChan = make(chan *BufferData)

	for {
		readyQueue <- w

		c := <-w.workerChan

		if c == nil {
			wg.Done()
			break
		}

		startPoint := 0
		for i := 0; i < c.n; i = i + 8 {
			p := binary.LittleEndian.Uint64(w.buffer[i : i+8])
			if i+7 > c.n {
				p = binary.LittleEndian.Uint64(w.refillBuffer(i, c.n))
			}

			if nl := findPosition(p, line_jump_mask); nl != -1 {
				w.saveLine(startPoint, i+nl)
				startPoint = i + nl + 1
			}
		}
	}
}

func (w *worker) saveLine(start int, end int) {
	for i := start; i < end; i = i + 8 {
		var p uint64
		if i+7 > end {
			p = binary.LittleEndian.Uint64(w.refillBuffer(i, end))
		} else {
			p = binary.LittleEndian.Uint64(w.buffer[i : i+8])
		}

		if psc := findPosition(p, semicolon_mask); psc != -1 {
			w.saveLineSeparated(start, i+psc, end)
			break
		}
	}
}

func (w *worker) saveLineSeparated(start int, midSeparator int, end int) {
	if midSeparator == -1 {
		return
	}
	temperature := parseTemperature(w.buffer[midSeparator+1 : end])

	h := hash(w.buffer[start:midSeparator])

	if v, ok := w.swdata.Get(h); ok {
		v.Count++
		v.Sum += temperature
		if v.Max < temperature {
			v.Max = temperature
		}
		if v.Min > temperature {
			v.Min = temperature
		}
	} else {
		w.swdata.Put(h, &stationData{Count: 1, Max: temperature, Min: temperature, Sum: temperature, Name: string(w.buffer[start:midSeparator])})
	}
}

func (w *worker) refillBuffer(i int, n int) []byte {
	b := make([]byte, 8)
	for j := 0; j < 8; j++ {
		if i+j <= n {
			b[j] = w.buffer[i+j]
		}
	}
	return b
}
