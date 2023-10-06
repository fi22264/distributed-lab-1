package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
	// TODO: all
	// Deal with an error event.
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	for {
		conn, e := ln.Accept()
		handleError(e)
		conns <- conn
	}
	// TODO: all

	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	reader := bufio.NewScanner(client)
	for reader.Scan() {
		msg_ := reader.Text()
		msg := Message{
			sender:  clientid,
			message: fmt.Sprintf("[%d]: %s\n", clientid, msg_),
		}

		msgs <- msg
	}
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.

	ln, e := net.Listen("tcp", *portPtr)
	if e != nil {
		panic(e)
	}

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			cID := len(clients) + 1
			fmt.Printf("New connection with ID: %d and addr: %s\n", cID, conn.RemoteAddr().String())
			clients[cID] = conn
			go handleClient(conn, cID, msgs)

		case msg := <-msgs:
			for i, client := range clients {
				if i == msg.sender {
					continue
				}
				writer := bufio.NewWriter(client)
				_, e := writer.WriteString(msg.message)
				handleError(e)
				e = writer.Flush()
				handleError(e)
			}
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
		}
	}
}
