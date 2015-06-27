package entropy

import (
	"io"
	"math"
	"os"
)

func NewEntropyCounter() *EntropyCounter {
	return &EntropyCounter{}
}

type EntropyCounter struct {
	counter  [256]int
	numBytes int
}

func (ec *EntropyCounter) AddByte(b byte) {
	ec.counter[int(b)]++
	ec.numBytes++
}

func (ec *EntropyCounter) ShannonEntropy() float64 {
	var entropy = float64(0)

	for _, cnt := range ec.counter {
		if cnt == 0 {
			continue
		}

		p := float64(cnt) / float64(ec.numBytes)

		entropy -= p * math.Log2(p)
	}

	return entropy
}

type Block interface {
	ShannonEntropy() float64
}

type block struct {
	data []byte
}

func (b *block) ShannonEntropy() float64 {

	counter := NewEntropyCounter()

	for _, b := range b.data {
		counter.AddByte(b)
	}

	return counter.ShannonEntropy()
}

func NewFileReaderFromFilename(filename string) (*FileReader, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	return NewFileReader(file), nil
}

func NewFileReader(rdr io.Reader) *FileReader {
	return &FileReader{rdr}
}

type FileReader struct {
	io.Reader
}

func (fr *FileReader) ShannonEntropy() float64 {
	var (
		buf     = make([]byte, 1024)
		counter = NewEntropyCounter()
		n       int
		err     error
	)

	for err == nil {
		n, err = fr.Read(buf)

		for i := 0; i < n; i++ {
			counter.AddByte(buf[i])
		}
	}

	return counter.ShannonEntropy()
}

type BlockReader struct {
	rdr io.Reader
	buf []byte
}

func NewBlockReaderFromFilename(filename string, blockSize int) (*BlockReader, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	return NewBlockReader(file, blockSize), nil
}

func NewBlockReader(rdr io.Reader, blockSize int) *BlockReader {
	return &BlockReader{
		rdr: rdr,
		buf: make([]byte, blockSize),
	}
}

func (er *BlockReader) Read() (Block, error) {
	n, err := er.rdr.Read(er.buf)

	return &block{
		data: er.buf[:n],
	}, err
}
