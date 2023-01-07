package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"golang.design/x/clipboard"
)

var once = sync.Once{}

func clipboardInit() {
	once.Do(func() {
		err := clipboard.Init()
		if err != nil {
			fmt.Fprintf(os.Stderr, "iit clipboard %+v", err)
			os.Exit(255)
		}
	})
}

type ClipboardReadCloser struct {
	data  []byte
	index int
}

func NewClipboardReadCloser() *ClipboardReadCloser {
	clipboardInit()
	data := clipboard.Read(clipboard.FmtText)
	return &ClipboardReadCloser{
		data: data,
	}
}

func (c *ClipboardReadCloser) Read(p []byte) (n int, err error) {
	if c.index >= len(c.data) {
		return 0, io.EOF
	}
	n = copy(p, c.data[c.index:])
	c.index += n
	return
}

func (c *ClipboardReadCloser) Close() error {
	c.data = c.data[0:0]
	c.index = 0
	return nil
}

type ClipboardWriteCloser struct {
	data []byte
}

func NewClipboardWriteCloser() *ClipboardWriteCloser {
	return &ClipboardWriteCloser{}
}

func (c *ClipboardWriteCloser) Write(p []byte) (n int, err error) {
	c.data = append(c.data, p...)
	return len(p), nil
}

func (c *ClipboardWriteCloser) Close() error {
	clipboard.Write(clipboard.FmtText, c.data)
	c.data = c.data[0:0]
	return nil
}
