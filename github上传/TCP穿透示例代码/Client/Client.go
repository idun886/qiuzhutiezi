package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"sync"
)

func FrowardTCP(fromConn, toConn net.Conn) {
	defer fromConn.Close()
	defer toConn.Close()
	go func() {
		_, err := io.Copy(toConn, fromConn)
		if err != nil {
			//todo log
			//fmt.Println("读取数据异常", err)
		}
	}()
	_, err := io.Copy(fromConn, toConn)
	if err != nil {
		//fmt.Println("读取数据异常", err)
		//todo log
	}

}

func HandConn(remoteConn net.Conn, localAddr string) {
	_, err := remoteConn.Write([]byte("OKKK"))
	if err != nil {
		remoteConn.Close()
		fmt.Printf("OKKK发送失败: %s\n", err)
		return
	}
	localConn, err := net.Dial("tcp", localAddr)
	if err != nil {
		remoteConn.Close()
		fmt.Printf("请求本地服务失败：%s\n", err)
		return
	}
	fmt.Printf("本地服务地址：%s\n", localConn.RemoteAddr().String())
	fmt.Printf("有人访问本地服务\n")
	go FrowardTCP(localConn, remoteConn)
}

func DealWith(LocalAddr, RemoteAddr string) {
	//与服务端建立连接
	RemoteConn, err := net.Dial("tcp", RemoteAddr)
	if err != nil {
		fmt.Printf("Dial remote addr error: %s\n", err)
		return
	}
	fmt.Printf("服务端连接成功 ->地址：%s\n", RemoteConn.RemoteAddr().String())
	//发送hello认证
	_, err = RemoteConn.Write([]byte("Hello"))
	if err != nil {
		RemoteConn.Close()
		fmt.Printf("Write remoteconn error: %s\n", err)
		return
	}
	signbuf := make([]byte, 64)
	//读取hi指令
	_, err = io.ReadFull(RemoteConn, signbuf[:2])
	if err != nil {
		RemoteConn.Close()
		return
	}
	//如果收到的信息和服务端发来的Hi回复不一样则提示
	if !bytes.Equal(signbuf[:2], []byte("Hi")) {
		fmt.Printf("鬼\n")
		return
	}
	fmt.Printf("okkkkk!!!!! 建立成功\n")
	for {
		_, err = io.ReadFull(RemoteConn, signbuf[:4])
		//如果收到服务端发开的CONN指令 代表收到一条请求连接信息
		if bytes.Equal(signbuf[:4], []byte("CONN")) {
			go HandConn(RemoteConn, LocalAddr)
			return
		}
	}

}
func main() {

	var confile string
	var wg sync.WaitGroup
	flag.StringVar(&confile, "c", "config.json", "指定配置文件")
	flag.Parse()
	conf := NewConfFileInfo(confile)
	conf, err := conf.ReadConfFile()
	if err != nil {
		fmt.Println(err)
	}
	confcontent, err := conf.ParserConf()
	if err != nil {
		fmt.Println(err)
		return
	}

	if confcontent != nil {
		wg.Add(len(confcontent.Client))
		for _, value := range confcontent.Client {
			go func(localServerAddr, remoteAddr string) {
				defer wg.Done()
				for {
					DealWith(localServerAddr, remoteAddr)
				}
			}(value.LocalServerAddr, value.RemoteAddr)
		}
	}
	wg.Wait()

}
