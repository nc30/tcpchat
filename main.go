package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
)

var (
	id    uuid.UUID
	count int = 0
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	defer log.Println("disconnected:", conn.RemoteAddr())
	reader := bufio.NewReaderSize(conn, 4098)

	log.Println("connect:", count, "from:", conn.RemoteAddr())
	fmt.Fprintln(conn, "HELLO", id.String())
	fmt.Fprintln(conn, "Your", count, "st connection.")

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(string(line))
		if string(line) == "add" {
			count++
			log.Println("countup", count)
			fmt.Fprintln(conn, "Your", count, "st connection.")
		}
		if string(line) == "bye" {
			break
		}
	}
}

func main() {
	id = uuid.New()
	log.SetFlags(0)
	log.SetPrefix(id.String() + " ")

	log.Println("listening:", "23")
	srv, err := net.Listen("tcp", ":23")
	if err != nil {
		panic(err)
	}

	for {
		c, err := srv.Accept()
		count++
		if err != nil {
			c.Close()
			return
		}

		go handleConn(c)
	}
}
