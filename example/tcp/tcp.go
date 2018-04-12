package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.1.107:1567")

	if err != nil {
		log.Println(err)
		return
	}

	go handleConn(conn)

	fmt.Fprintf(conn, "hahah")

}

func handleConn(conn net.Conn) {
	for {
		var data []byte
		n, _ := conn.Read(data)

		fmt.Println(string(data[:n]))
	}

}
