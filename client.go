package chat

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"time"
)

// Client will be used to send and recevie the message
type Client struct {
	name       string
	encoder    *gob.Encoder
	decoder    *gob.Decoder
	connection net.Conn
	buf        bytes.Buffer
}

// NewClient is created based on the name given
func NewClient(name string) {
	var connection net.Conn
	var err error
	var done chan bool = make(chan bool)
	gob.Register(net.TCPAddr{})
	if connection, err = net.Dial("tcp", "localhost:8080"); err != nil {
		panic("could not connect to server")
	}
	var c *Client = &Client{
		name:       name,
		connection: connection,
	}
	c.encoder = gob.NewEncoder(connection)
	c.decoder = gob.NewDecoder(connection)
	var inp chan string = make(chan string, 1)
	go c.Cmd(inp, done)
	for {
		select {
		case d := <-inp:
			c.Send(d, nil)
		case <-done:
			connection.Close()
			os.Exit(0)
			break
		}
	}
}

// Send data to a client
func (c *Client) Send(data string, to net.Addr) {
	var packet *Message = &Message{
		ContentType: "text/plain",
		From:        c.connection.LocalAddr().String(),
		Data:        data,
	}
	if err := c.encoder.Encode(packet); err != nil {
		fmt.Println("error: encoding packet", err)
	}
}

// Cmd interface for the client
func (c *Client) Cmd(in chan<- string, done chan<- bool) {
	var scanner = bufio.NewScanner(os.Stdin)
	for {
		var data string
		fmt.Print(c.name, "->")
		if scanner.Scan() {
			data = scanner.Text()
		}
		if data == "quit" {
			c.connection.Close()
			done <- true
			break
		}
		in <- data
		time.Sleep(1 * time.Millisecond)
	}
}

// Receive data from other clients
func (c *Client) Receive() {
	var data bytes.Buffer
	if _, err := c.connection.Read(data.Bytes()); err != nil {
		fmt.Println("error: reading data from pipe")
	}
}
