package srv

import (
	"bufio"
	"bytes"
	"context"
	"io"
)

type Request struct {
	Reader  bufio.Reader
	Context context.Context
	Data    interface{}
}

func (r *Request) Read() ([]byte, error) {
	var buffer bytes.Buffer
	for {
		b, prefix, err := r.Reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		buffer.Write(b)
		if !prefix {
			break
		}
	}

	return buffer.Bytes(), nil
}
