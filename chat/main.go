package main

import (
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:7007")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer listener.Close()

	connection, err := listener.Accept()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer connection.Close()

	io.Copy(connection, connection)
}
