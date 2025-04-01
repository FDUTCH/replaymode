package format

import (
	"bufio"
	"os"
	"strings"
)

const suffix = ".mcreplay"

func Create(name string) (ByteWriter, error) {
	if !strings.HasSuffix(name, suffix) {
		name += suffix
	}
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	return &byteWriter{f: f, Writer: bufio.NewWriter(f)}, err
}

type byteWriter struct {
	f *os.File
	*bufio.Writer
}

func (b *byteWriter) Close() error {
	b.Flush()
	return b.f.Close()
}

func Open(name string) (ByteReader, error) {
	if !strings.HasSuffix(name, suffix) {
		name += suffix
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return &byteReader{File: f, b: make([]byte, 1)}, nil
}

type byteReader struct {
	*os.File
	b []byte
}

func (b *byteReader) ReadByte() (byte, error) {
	_, err := b.Read(b.b)
	return b.b[0], err
}
