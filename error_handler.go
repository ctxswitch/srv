package srv

import "net"

type ErrorHandler interface {
	Send(conn net.Conn)
}
