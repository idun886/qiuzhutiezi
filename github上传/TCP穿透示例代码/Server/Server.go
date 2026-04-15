package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"sync"
)

func HandleAccpet(openListener net.Listener, openConnChan chan net.Conn) {
	for {
		openconn, err := openListener.Accept()
		if err != nil {
			fmt.Printf("openListener accept error :%s\n", err)
		}
		openConnChan <- openconn
	}
}

func HandleConn(tunnelConn net.Conn, openConnChan chan net.Conn) {
	signbuf := make([]byte, 64)
	//Hello
	buflen, err := io.ReadFull(tunnelConn, signbuf[:5])
	if err != nil {
		fmt.Printf("read signbuf error :%s\n", err)
		tunnelConn.Close()
		return
	}
	if !bytes.Equal(signbuf[:5], []byte("Hello")) {
		fmt.Printf("client can't auth,client sign:%s\n", string(signbuf[:buflen]))
		tunnelConn.Close()
		return
	}
	fmt.Printf("收到client sign:%s\n", string(signbuf[:5]))
	_, err = tunnelConn.Write([]byte("Hi"))
	if err != nil {
		fmt.Printf("tunnelConn write error:%s\n", err)
		tunnelConn.Close()
		return
	}
	fmt.Printf("发送 Hi sign 成功\n")
	for {
		select {
		case openconn := <-openConnChan:
			buf := []byte("CONN")
			_, err = tunnelConn.Write(buf)
			if err != nil {
				fmt.Printf("tunnelConn write error :%s\n", err)
				tunnelConn.Close()
				openConnChan <- openconn
				return
			}
			for {
				_, err = io.ReadFull(tunnelConn, buf[:4])
				if err != nil {
					fmt.Printf("tunnelConn readFull error:%s\n", err)
					tunnelConn.Close()
					openConnChan <- openconn
					return
				}
				if bytes.Equal(buf[:4], []byte("OKKK")) {
					break
				}
			}
			go FrowardTCP(tunnelConn, openconn)
			return
		}
	}

}

func FrowardTCP(fromConn, toConn net.Conn) {
	defer fromConn.Close()
	defer toConn.Close()
	go func() {
		_, err := io.Copy(toConn, fromConn)
		if err != nil {
			//todo log
		}
	}()
	_, err := io.Copy(fromConn, toConn)
	if err != nil {
		//todo log
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
		wg.Add(len(confcontent.Server))
		for _, value := range confcontent.Server {
			go func(openport, tunnelport string) {
				defer wg.Done()

				TunnelListener, err := net.Listen("tcp", "0.0.0.0:"+tunnelport)
				if err != nil {
					fmt.Printf("本地端口监听异常，或被占用：%s\n", err)
					return
				}
				OpenListener, err := net.Listen("tcp", "0.0.0.0:"+openport)
				if err != nil {
					fmt.Printf("监听外部访问端口异常，或被占用：%s\n", err)
					return
				}
				//创建连接管道
				openconnchan := make(chan net.Conn)
				go HandleAccpet(OpenListener, openconnchan)

				for {
					TunnelConn, err := TunnelListener.Accept()
					if err != nil {
						fmt.Printf("Tunnel Accept error :%s\n", err)
					}
					go HandleConn(TunnelConn, openconnchan)
				}

			}(value.OpenPort, value.TunnelPort)
		}
	}
	wg.Wait()

}
