package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func tcp_client() {
	fmt.Println("Connecting")
	conn, err := net.Dial("tcp", "localhost:8080")
	fmt.Println("Connected")
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(status)
	conn.Write([]byte("Message to the server\n"))
}

func tcp_server() {
	go func() {
		ln, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatal(err)
		}
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
			}
			conn.Write([]byte("Testing out\n"))
			fmt.Println(conn.Read([]byte{}))
			// status, err := bufio.NewReader(conn).ReadString('\n')
			// fmt.Println(status)
		}
	}()
	time.Sleep(time.Millisecond * 100)
}

func test_tcp() {
	tcp_server()
	tcp_client()
}
