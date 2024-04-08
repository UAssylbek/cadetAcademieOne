package main

import (
	"fmt"
	"net"
	netcat "netcat/app"
	"os"
)

func main() {
	var port string
	if len(os.Args) == 1 {
		port = ":8989"
	} else if len(os.Args) == 2 {
		port = ":" + os.Args[1]
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	// Старт сервера
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Listening on the port", port)

	// Обработка входящих подключений
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go netcat.HandleConnection(conn)
	}
}
