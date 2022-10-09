package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	//wg := sync.WaitGroup{}
	//wg.Add(100)
	//for i := 0; i < 100; i++ {
	//	go func(i int) {
	//		fmt.Println(i)
	//		wg.Done()
	//	}(i)
	//}
	//wg.Wait()
	//fmt.Println("123")

	//listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), port, ""})
	//if err != nil {
	//	log.Println("监听端口失败:", err.Error())
	//	return
	//}
	//log.Println("已初始化连接，等待客户端连接...")
	//Server(listen)

	go startServer("server1", ":8838")
	go startServer("server2", ":8848")

	go startClient("client1", 3, ":8848")
	startClient("client2", 2, ":8848")
}

func startServer(name string, port string) {
	fmt.Println("Starting Tcp Server..., name: " + name)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err.Error())
	}
	go func() {
		defer func(listener net.Listener) {
			err := listener.Close()
			if err != nil {
				fmt.Println(err.Error())
			}
		}(listener)
		for {
			server, err := listener.Accept()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			go handlerServer1(server, name)
		}
	}()
}

func startClient(name string, _time time.Duration, port string) {
	fmt.Println("Starting Tcp Client...,name: " + name)
	client, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for {
		handlerClient(client, name, _time)
	}
}

func handlerServer(conn net.Conn, name string) {
	pc, err := net.Dial("tcp", ":8838")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer func(pc net.Conn) {
		err := pc.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(pc)
	for {
		buf := make([]byte, 1024)
		_len, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		_input := string(buf[:_len])
		var _result string
		if strings.Contains(_input, "ping") {
			if name == "server2" {
				_, err := pc.Write([]byte(_input))
				if err != nil {
					fmt.Println(err.Error())
				}
			} else {
				_result = _input
				fmt.Println("hello from " + name + ", " + _result)
			}
		} else {
			_result = _input
			fmt.Println("hello from " + name + ", " + _result)
		}
	}
}

func handlerServer1(conn net.Conn, name string) {
	if strings.Contains(name, "server2") {
		forward(conn)
	} else {
		buf := make([]byte, 1024)
		for {
			writeLen(conn, buf)
		}
	}
}

func writeLen(conn net.Conn, buf []byte) string {
	_len, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err.Error())
	}
	_input := string(buf[:_len])
	var _result string
	_result = _input
	fmt.Println("hello from server1," + _result)
	return _result
}

func forward(conn net.Conn) {
	tonn, err := net.Dial("tcp", ":8838")
	if err != nil {
		fmt.Println(err.Error())
	}
	var wg sync.WaitGroup
	go func(conn, tonn net.Conn) {
		wg.Add(1)
		defer wg.Done()
		if _, err := io.Copy(conn, tonn); err != nil {
			fmt.Println(err.Error())
		}
		conn.Close()
	}(conn, tonn)
	wg.Add(1)
	defer wg.Done()
	if _, err := io.Copy(tonn, conn); err != nil {
		fmt.Println(err.Error())
	}
	tonn.Close()
	wg.Wait()
}

func handlerClient(conn net.Conn, name string, _time time.Duration) {
	for {
		_timestamp := time.Now().UnixNano()
		var _input string
		if name == "client1" {
			_input = "hello from " + name + ", ping"
		} else {
			_input = "hello from " + name + ", " + strconv.FormatInt(_timestamp, 10)
		}
		_, err := conn.Write([]byte(_input))
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(time.Second * _time)
	}
}
