package main

import (
	"context"
	"log"
	"net"
	"strings"
)

func readConn(ctx context.Context, conn net.Conn, message chan<- string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data := make([]byte, 1024)
			_, err := conn.Read(data)
			if err != nil {
				log.Printf("读取错误: %v\n", err)
				return
			}

			strData := string(data)
			log.Println("Received:", strData)

			message <- strData
		}
	}
}

func writeConn(ctx context.Context, conn net.Conn, message <-chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			strData := <-message
			_, err := conn.Write([]byte(strData))
			if err != nil {
				log.Printf("写入错误: %v\n", err)
				return
			}

			log.Println("Send:", strData)
		}
	}
}

func doServerStuff(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	readChan := make(chan string)
	writeChan := make(chan string)

	go readConn(ctx, conn, readChan)
	go writeConn(ctx, conn, writeChan)

	for {
		select {
		case readStr := <-readChan:
			upper := strings.ToUpper(readStr)
			writeChan <- upper
		case <-ctx.Done():
			break
		}
	}

}

func main() {
	listener, err := net.Listen("tcp", ":9900")
	if err != nil {
		log.Fatalf("服务启动错误: %v", err)
	}
	defer listener.Close()

	log.Println("Starting the server ...", listener.Addr().String())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go doServerStuff(ctx, conn)
	}
}
