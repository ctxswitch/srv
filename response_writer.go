package srv

import (
	"bufio"
	"net"
)

type ResponseWriter interface {
	Write([]byte) (int, error)
}

type Writer struct {
	Connection net.Conn
	Data       interface{}
}

func (w Writer) Write(b []byte) (n int, err error) {
	writer := bufio.NewWriter(w.Connection)
	n, err = writer.WriteString(string(b))
	if err == nil {
		err = writer.Flush()
	}
	return
}
