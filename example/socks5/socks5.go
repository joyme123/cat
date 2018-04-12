package main

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/proxy"
)

func main() {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, nil)

	if err != nil {
		log.Println(err)
		return
	}

	conn, err := dialer.Dial("tcp", "47.98.115.70:10007")

	if err != nil {
		log.Println(err)
		return
	}

	go handleConn(conn)

	fmt.Fprint(conn, "dasdasd")

	for {

	}

}

func handleConn(conn net.Conn) {
	for {
		var data [1024]byte
		n, err := conn.Read(data[:])
		if err != nil {
			log.Println(err)
			return
		}

		if n > 0 {
			fmt.Println(string(data[:n]))
		}
	}
}
