package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net"
	"os"
)

func encode(b []byte) string {
	return base64.StdEncoding.Strict().EncodeToString(b)
}

func encrypt(text, secret string) string {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		println("unable to encrypt", err.Error())
		os.Exit(1)
	}

	plainText := []byte(text)

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		println("unable to create gcm", err.Error())
		os.Exit(1)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		println("unable to create nonce", err.Error())
		os.Exit(1)
	}

	return encode(gcm.Seal(nonce, nonce, plainText, nil))
}

func main() {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "supersecretpasswordatleast32byte"
	}

	str := encrypt(os.Args[1], secret) + "\r\n"

	raddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9000")
	if err != nil {
		println("resolve failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		println("dial failed:", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte(str))
	if err != nil {
		println("write failed:", err.Error())
		os.Exit(1)
	}

	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		println("read failed:", err.Error())
		os.Exit(1)
	}

	println("[REPLY]:", string(reply))

	conn.Close()
}
