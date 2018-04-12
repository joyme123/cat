package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"

	"golang.org/x/net/proxy"
)

func main() {

	host := flag.String("h", "127.0.0.1:8900", "http代理监听的ip地址")
	socksAddr := flag.String("s", "127.0.0.1:1080", "socks的连接ip和端口")
	flag.Parse()

	proxy, err := net.Listen("tcp", *host)
	defer proxy.Close()
	if err != nil {
		log.Println(err)
	}

	for {
		client, err := proxy.Accept()
		if err != nil {
			log.Println(err)
		}

		go handleHTTPProxy(client, *socksAddr)
	}

}

func handleHTTPProxy(cli net.Conn, socksAddr string) {
	defer cli.Close()
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

	dialer, err := proxy.SOCKS5("tcp", socksAddr, nil, nil) //开启socks5连接

	if err != nil {
		log.Printf("socks5连接出错%v", err)
		return
	}
	target, err := dialer.Dial("tcp", address) // socks5对请求进行拨号
	defer target.Close()
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
