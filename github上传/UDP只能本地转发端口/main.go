package main

import (
	"fmt"
	"net"
)

var clientMap = make(map[string]*net.UDPAddr)

func main() {
	localudpaddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:5889")
	openudpaddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9999")
	openconn, err := net.ListenUDP("udp", openudpaddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 用于和目标通信
	localconn, err1 := net.DialUDP("udp", nil, localudpaddr)
	if err != nil {
		panic(err1)
	}
	defer localconn.Close()
	fmt.Println("UDP forward: :9999 -> 127.0.0.1:5889")
	go func() {
		buf := make([]byte, 2048)
		for {
			num, _, err3 := localconn.ReadFromUDP(buf)
			if err3 != nil {
				fmt.Println(err)
				continue
			}
			for _, clientAddr := range clientMap {
				openconn.WriteToUDP(buf[:num], clientAddr)
			}

		}

	}()

	buf := make([]byte, 2048)
	for {
		num, addr, err2 := openconn.ReadFromUDP(buf)
		if err2 != nil {
			fmt.Println(err2)
			continue
		}
		clientMap[addr.String()] = addr
		// 转发给目标
		_, err = localconn.Write(buf[:num])
	}
}
