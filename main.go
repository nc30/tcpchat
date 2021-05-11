package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/google/uuid"
)

var (
	id          uuid.UUID
	count       int = 0
	addr        string
	mux         sync.Mutex
	connections map[net.Conn]bool
)

func init() {
	connections = map[net.Conn]bool{}
}

func broadcast(data ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	for c, _ := range connections {
		fmt.Fprintln(c, data...)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	defer log.Println("disconnected:", conn.RemoteAddr())
	mux.Lock()
	connections[conn] = true
	mux.Unlock()
	defer func() {
		mux.Lock()
		delete(connections, conn)
		mux.Unlock()
	}()

	reader := bufio.NewReaderSize(conn, 4098)

	log.Println("connect:", count, "from:", conn.RemoteAddr())
	fmt.Fprintln(conn, "HELLO", id.String())
	fmt.Fprintln(conn, "Now", len(connections), "connections alive.")

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(string(line))

		commands := strings.Split(string(line), " ")
		var cmd string
		if len(commands) != 0 {
			cmd = strings.ToLower(commands[0])
		} else {
			cmd = strings.ToLower(string(line))
		}

		if cmd == "add" {
			count++
			log.Println("countup", count)
			broadcast("countup to", count)
		}
		if cmd == "sub" {
			count--
			log.Println("countdown", count)
			broadcast("countup to", count)
		}
		if cmd == "bye" || cmd == "quit" || cmd == "exit" {
			break
		}
	}
}

func main() {

	id = uuid.New()
	log.SetFlags(0)
	log.SetPrefix(id.String() + " ")

	flag.StringVar(&addr, "l", ":23", "listen address")
	flag.Parse()

	log.Println("listening:", addr)
	srv, err := net.Listen("tcp", addr)
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
