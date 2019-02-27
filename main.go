package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

type csvStreamer struct {
	enc *csv.Reader
	w   io.Writer
}

func newCSVStreamer() *csvStreamer {
	pr, pw := io.Pipe()
	return &csvStreamer{
		enc: csv.NewReader(pr),
		w:   pw,
	}
}

type syncResp struct {
	rows [][]string
	err  error
}

type writerTo interface {
	WriteTo(w io.Writer) (n int64, err error)
}

func (csv *csvStreamer) Sync(conn writerTo) (rows [][]string, err error) {
	go func() {
		conn.WriteTo(csv.w)
	}()

	for {
		row, err := csv.enc.Read()
		if err != nil {
			return nil, err
		} else if row == nil {
			break
		}
		rows = append(rows, row)
	}

	return rows, nil
}

func main() {
	rawTestData := "" +
		"123,abc,Organic\n" +
		"456,def,Organic\n" +
		"789,ghi,Organic\n"

	buf := bytes.NewBufferString(rawTestData)

	stream := newCSVStreamer()

	rows, err := stream.Sync(buf)
	fmt.Printf("(main) rows %v, err %v\n", rows, err)
	time.Sleep(time.Second)
}
