package format

import (
	"encoding/binary"
	"io"
	"time"
)

func NewTimeReader(reader io.Reader) func() (time.Duration, error) {

	buff := make([]byte, 4)

	return func() (time.Duration, error) {
		_, err := reader.Read(buff)
		if err != nil {
			return 0, err
		}

		return time.Duration(binary.LittleEndian.Uint32(buff)) * time.Millisecond, nil
	}
}

func NewTimeWriter(writer io.Writer) func() error {

	buff := make([]byte, 4)

	return func() error {
		stamp := time.Now().UnixMilli()

		binary.LittleEndian.PutUint32(buff, uint32(stamp))

		_, err := writer.Write(buff)
		return err
	}
}
