package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ip   = ""
	port = 8899
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

	startServer("server1", ":8838")
	startServer("server2", ":8848")

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
		defer listener.Close()
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
	defer pc.Close()
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
	defer conn.Close()
	if strings.Contains(name, "server2") {
		pc, err := net.Dial("tcp", ":8838")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer pc.Close()
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err1 := io.Copy(pc, conn)
			if err1 != nil {
				fmt.Println(err1.Error())
				return
			}
		}()
		wg.Wait()
		_, err2 := io.Copy(conn, pc)
		if err2 != nil {
			fmt.Println(err2.Error())
			return
		}
	} else {
		WriteLen(conn)
	}
}

func WriteLen(conn net.Conn) string {
	buf := make([]byte, 1024)
	_len, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err.Error())
	}
	_input := string(buf[:_len])
	var _result string
	_result = _input
	return _result
	//fmt.Println("hello from ," + _result)
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

func Server(listen *net.TCPListener) {
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("接受客户端连接异常:", err.Error())
			continue
		}
		log.Println("客户端连接来自:", conn.RemoteAddr().String())
		defer conn.Close()
		go func() {
			data := make([]byte, 8192)
			for {
				i, err := conn.Read(data)
				log.Println("客户端发来数据:", i)
				if err != nil {
					log.Println("读取客户端数据错误:", err.Error())
					break
				}
				//go conn.Write(Send(data))
				go func() {
					dealData(conn)
				}()
			}
		}()
	}
}

func dealData(conn *net.TCPConn) {
	pTCPConn, err := net.Dial("tcp", "127.0.0.1:80")
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
	defer pTCPConn.Close()
	go func() {
		io.Copy(pTCPConn, conn)
	}()
	io.Copy(conn, pTCPConn)
}

func Send(data []byte) (buf []byte) {
	pTCPConn, err := net.Dial("tcp", "127.0.0.1:80")
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
	n, errWrite := pTCPConn.Write(data)
	if errWrite != nil {
		log.Printf("Error: %s", errWrite.Error())
		return
	}
	defer pTCPConn.Close()
	log.Printf("writed: %d\n", n)
	buf, errRead := ioutil.ReadAll(pTCPConn)
	log.Println("服务端发来数据:", len(buf))
	if errRead != nil {
		log.Printf("Error: %s", errRead.Error())
		return
	}
	return
}
