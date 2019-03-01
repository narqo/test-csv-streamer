package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	flag.Parse()
	filePath := flag.Arg(0)
	if filePath == "" {
		log.Fatal("path is empty")
	}

	if err := run(filePath); err != nil {
		log.Fatal(err)
	}
}

type stream struct {
	dec *csv.Reader
	w   *io.PipeWriter
}

func newStream() *stream {
	pr, pw := io.Pipe()
	return &stream{
		dec: csv.NewReader(pr),
		w:   pw,
	}
}

func (s *stream) Write(p []byte) (n int, err error) {
	return s.w.Write(p)
}

func (s *stream) Close() error {
	return s.w.Close()
}

func (s *stream) NextRecord() (rec []string, err error) {
	return s.dec.Read()
}

func run(filePath string) (err error) {
	stream := newStream()

	go writeTo(stream)

	for err == nil {
		var rec []string
		rec, err = stream.NextRecord()
		if rec != nil {
			fmt.Fprintf(os.Stdout, "record: %v\n", rec)
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}

// emulate something that accepts only io.Writer
// e.g. a net driver that implements io.WriterTo
func writeTo(sw io.WriteCloser) {
	rawcsv := []byte("a,b,c,,123,456\nx,y,z,0,333,a\n")

	// write a byte at a time, emulating slow source
	data := rawcsv
	for ; len(data) > 0; data = data[1:] {
		sw.Write(data[:1])
	}

	data = rawcsv
	for ; len(data) > 0; data = data[1:] {
		sw.Write(data[:1])
	}

	sw.Close()
}
