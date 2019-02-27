package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type csvstream struct {
	buf *bytes.Buffer
	csv *csv.Reader

	out io.Writer
}

func newStream(out io.Writer) *csvstream {
	buf := bytes.NewBuffer(nil)
	return &csvstream{
		buf: buf,
		csv: csv.NewReader(buf),
		out: out,
	}
}

func (s *csvstream) Write(p []byte) (n int, err error) {
	n, err = s.buf.Write(p)
	if err != nil {
		return 0, err
	}

	for {
		rec, err := s.csv.Read()
		if err != nil {
			break
		}
		// map records and pass them to out writer
		fmt.Fprintf(s.out, "record: %v\n", rec)
	}
	if err == io.EOF {
		err = nil
	}

	return n, err
}

func main() {
	flag.Parse()
	filePath := flag.Arg(0)
	if filePath == "" {
		log.Fatal("path is empty")
	}

	if err := run(filePath); err != nil {
		panic(err)
	}
}

func run(filePath string) error {
	s := newStream(os.Stderr)
	return streamTo(s, filePath)
}

// <some function> which receives writer and streams csv to it
func streamTo(w io.Writer, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	return err
}
