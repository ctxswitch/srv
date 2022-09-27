package main

import (
	"bytes"
	"fmt"

	"ctx.sh/srv"
)

type EchoHandler struct{}

func (h EchoHandler) ServeTCP(w srv.ResponseWriter, r *srv.Request) {
	fmt.Println("[EchoHandler]: entry")
	_, err := w.Write(r.Data.([]byte))
	if err != nil {
		return
	}
	fmt.Println("[EchoHandler]: exit")
}

type Transformer struct{}

func (t *Transformer) Middleware(next srv.Handler) srv.Handler {
	return srv.HandlerFunc(func(w srv.ResponseWriter, r *srv.Request) {
		fmt.Println("[Transformer]: entry")
		r.Data = bytes.ToUpper(r.Data.([]byte))

		next.ServeTCP(w, r)
		fmt.Println("[Transformer]: exit")
	})
}

func Parser(next srv.Handler) srv.Handler {
	return srv.HandlerFunc(func(w srv.ResponseWriter, r *srv.Request) {
		fmt.Println("[Parse]: entry")

		b, err := r.Read()
		if err != nil {
			w.Write([]byte("error\r\n"))
			return
		}

		r.Data = b

		next.ServeTCP(w, r)
		fmt.Println("[Parse]: exit")
	})
}

func main() {
	t := Transformer{}

	rts := srv.NewRouter()
	rts.Handle(&EchoHandler{})
	rts.Use(Parser)
	rts.Use(t.Middleware)

	err := srv.ListenAndServe("tcp", "127.0.0.1:9000", rts)
	fmt.Println(err.Error())
}
