package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func read(conn net.Conn) {
	reader := bufio.NewScanner(conn)
	for reader.Scan() {
		msg := reader.Text()
		fmt.Println(msg)
	}
}

func write(conn net.Conn) {
	scan := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(conn)
	fmt.Println("Enter text")
	for scan.Scan() {
		text := scan.Text()

		_, e := writer.WriteString(text + "\n")
		if e != nil {
			panic(e)
		}
		e = writer.Flush()
		if e != nil {
			panic(e)
		}
		fmt.Println("[You]:", text)
	}
}

func handleConnection(conn net.Conn) {
	go read(conn)
	write(conn)
}

func main() {
	// Get the server address and port from the commandline arguments.
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()

	conn, e := net.Dial("tcp", *addrPtr)
	if e != nil {
		panic(e)
	}

	handleConnection(conn)
}
