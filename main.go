package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type csvstream struct {
	buf *bytes.Buffer
	csv *csv.Reader
}

func NewStream(out io.Writer) *csvstream {
	buf := bytes.NewBuffer(nil)
	return &csvstream{
		buf: buf,
		csv: csv.NewReader(buf),
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
		fmt.Printf("record: %v\n", rec)
	}
	if err == io.EOF {
		err = nil
	}

	return n, err
}

func main() {
	rawTestData := "" +
		"123,abc,Organic\n" +
		"456,def,Organic\n" +
		"789,ghi,Organic\n"

	buf := bytes.NewBufferString(rawTestData)

	s := NewStream(os.Stderr)
	_, err := buf.WriteTo(s)
	fmt.Println(err)
}
