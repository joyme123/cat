package main

import (
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:10007")
	defer listener.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			return
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	io.Copy(conn, conn)
}
