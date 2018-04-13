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

var pac []string

func main() {

	host := flag.String("h", "127.0.0.1:8900", "http代理监听的ip地址")
	socksAddr := flag.String("s", "127.0.0.1:1080", "socks的连接ip和端口")
	pacfile := flag.String("c", "/etc/cat/pac.txt", "pac文件地址")
	flag.Parse()

	log.Printf(*pacfile)

	//解析pac
	pac = parsePac(*pacfile)

	if len(pac) > 0 {
		log.Printf("解析到%d条url", len(pac))
	}

	proxy, err := net.Listen("tcp", *host)

	if err != nil {
		log.Println(err)
	}

	defer proxy.Close()

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
	var method, host, address, targetURL string
	fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
	hostPort, err := url.Parse(host)
	if err != nil {
		log.Fatal(err)
		return
	}

	if hostPort.Opaque == "443" { //https访问
		address = hostPort.Scheme + ":443"
		targetURL = hostPort.Scheme
	} else { //http访问
		if strings.Index(hostPort.Host, ":") == -1 { //不带端口，则默认80
			address = hostPort.Host + ":80"
		} else {
			address = hostPort.Host
		}
		targetURL = hostPort.Host
	}

	var target net.Conn
	var neterr error

	// 检查hostPort.Host是否在pac中
	if isInPac(targetURL, pac) {
		//走socks5代理
		dialer, err := proxy.SOCKS5("tcp", socksAddr, nil, nil) //开启socks5连接

		if err != nil {
			log.Printf("socks5连接出错%v", err)
			return
		}
		target, neterr = dialer.Dial("tcp", address) // socks5对请求进行拨号
		defer target.Close()
		if neterr != nil {
			log.Fatal(err)
			return
		}
	} else {
		//不走socks5代理
		target, err = net.Dial("tcp", address)

		defer target.Close()
		if err != nil {
			log.Fatal(err)
			return
		}

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
