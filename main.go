package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
)

type csvStreamer struct {
	w   io.Writer
	enc *csv.Reader
}

func newCSVStreamer() *csvStreamer {
	rw, wr := io.Pipe()
	return &csvStreamer{
		enc: csv.NewReader(rw),
		w:   wr,
	}
}

func (csv *csvStreamer) Write(p []byte) (int, error) {
	return csv.w.Write(p)
}

type syncResp struct {
	rows []string
	err  error
}

func (csv *csvStreamer) Sync(buf *bytes.Buffer) (rows []string, err error) {
	resp := make(chan syncResp, 1)

	go func() {
		rows, err := csv.enc.Read()
		resp <- syncResp{rows, err}
	}()

	_, err = io.Copy(csv, buf)
	if err != nil {
		return nil, err
	}

	r := <-resp
	return r.rows, r.err
}

func main() {
	rawTestData := "" +
		"123,abc,Organic\n" +
		"456,def,Organic\n" +
		"789,ghi,Organic\n"

	buf := bytes.NewBufferString(rawTestData)

	stream := newCSVStreamer()

	for i := 0; i < 3; i++ {
		rows, err := stream.Sync(buf)
		fmt.Printf("rows %v, err %v\n", rows, err)
	}
}
