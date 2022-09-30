package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"ctx.sh/srv"
)

type EchoHandler struct{}

func (h *EchoHandler) ServeTCP(w srv.ResponseWriter, r *srv.Request) {
	fmt.Println("[EchoHandler]: entry")
	_, err := w.Write(r.Data.([]byte))
	if err != nil {
		return
	}
	fmt.Println("[EchoHandler]: exit")
}

func TransformHandler(next srv.Handler) srv.Handler {
	return srv.HandlerFunc(func(w srv.ResponseWriter, r *srv.Request) {
		fmt.Println("[Transform]: entry")
		r.Data = bytes.ToUpper(r.Data.([]byte))

		next.ServeTCP(w, r)
		fmt.Println("[Transform]: exit")
	})
}

type Encryption struct {
	secret string
}

func (d *Encryption) decode(b []byte) ([]byte, error) {
	str, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

func (e *Encryption) decrypt(b []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(e.secret))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(b) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := b[:nonceSize], b[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func (e *Encryption) Handler(next srv.Handler) srv.Handler {
	return srv.HandlerFunc(func(w srv.ResponseWriter, r *srv.Request) {
		fmt.Println("[Encryption]: decrypt")
		// do I need to strip the newline?
		decoded, err := e.decode(r.Data.([]byte))
		if err != nil {
			r.Data = []byte("error")
			return
		}
		msg, err := e.decrypt(decoded)
		if err != nil {
			r.Data = []byte("error")
			return
		}
		r.Data = msg
		next.ServeTCP(w, r)
		// encrypt for write?
		fmt.Println("[Encryption]: exit")
	})
}

func ReadHandler(next srv.Handler) srv.Handler {
	return srv.HandlerFunc(func(w srv.ResponseWriter, r *srv.Request) {
		fmt.Println("[ReadHandler]: entry")

		b, err := r.Read()
		if err != nil {
			w.Write([]byte("error\r\n"))
			return
		}

		r.Data = b

		next.ServeTCP(w, r)
		fmt.Println("[ReadHandler]: exit")
	})
}

func main() {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "supersecretpasswordatleast32byte"
	}

	server := srv.Server{}

	rts := srv.NewRouter()
	rts.Handle(&EchoHandler{})

	rts.Use(ReadHandler)

	encryption := &Encryption{secret: secret}
	rts.Use(encryption.Handler)

	rts.Use(TransformHandler)

	err := server.ListenAndServe("tcp", "127.0.0.1:9000", rts)
	fmt.Println(err.Error())
}
