package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

func main() {

	listener, err := net.Listen("tcp", "127.0.0.1:8900")
	defer listener.Close()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("server start:127.0.0.1:8900")

	for {
		cli, err := listener.Accept()
		defer cli.Close()

		if err != nil {
			log.Fatal(err)
			break
		}
		go handleClientConn(cli)
	}
}

func handleClientConn(cli net.Conn) {
	var b [1024]byte

	n, err := cli.Read(b[:])

	if err != nil {
		log.Fatal(err)
		return
	}

	// if n > 0 {
	// 	fmt.Println(string(bytes[:]))
	// }
	var method, host, address string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	hostPort, err := url.Parse(host)
	if err != nil {
		log.Fatal(err)
		return
	}

	if hostPort.Opaque == "443" { //https访问
		address = hostPort.Scheme + ":443"
	} else { //http访问
		if strings.Index(hostPort.Host, ":") == -1 { //不带端口，则默认80
			address = hostPort.Host + ":80"
		} else {
			address = hostPort.Host
		}
	}

	target, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
		return
	}

	if method == "CONNECT" {
		fmt.Fprint(cli, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		target.Write(b[:n])
	}

	//进行转发
	go io.Copy(target, cli)
	io.Copy(cli, target)
}
