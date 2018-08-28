# cat

sock5代理转http代理，可用于linux环境。
windows需自己手动编译源代码

使用方式

## 默认

```
./cat
```
此时默认的http代理地址是127.0.0.1:8900，默认连接的socks5代理地址是127.0.0.1:1080

## 指定ip和端口号

```
./cat -h 127.0.0.1:8989 -s 127.0.0.1:1080
```

-h代表http代理的地址和端口
-s是已经存在的socks5代理地址和端口（比如shadowsocks）

之后在控制台中
```
export http_proxy=127.0.0.1:8989
export https_proxy=127.0.0.1:8989
```
## pac文件

程序默认加载的pac文件路径是/etc/cat/pac.txt
可以通过-c来指定pac地址

pac文件中的url每行代表一个，可以写成正则表达式。比如为了支持google.com的所有域名。可以写成(.*)google.com

# 编译

```
go build .
```