package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	//打开连接:
	conn, err := net.Dial("tcp", ":9900")

	if err != nil {
		fmt.Printf("error %v connecting!", err)
		os.Exit(1)
	}

	conn.Write([]byte("我是谁?\nA\nB\nC\nD\nE\nF"))
}
