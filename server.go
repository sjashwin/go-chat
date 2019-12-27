package chat

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
)

// Server manges the communications between two parties
type Server interface {
	Start()
	Stop()
	Commit()
}

// Clan is a name given to the server
type Clan struct {
	database Database
	listener net.Listener
	handle   func(conn net.Conn)
	route    map[string]net.Conn
}

// Start a server to monitor the communcation
func (c *Clan) Start() {
	var err error

	if c.listener, err = net.Listen("tcp", ":8080"); err != nil {
		panic("error: listening to port 8080 on TCP")
	}
	for {
		if connection, err := c.listener.Accept(); err != nil {
			panic("error: accepting incoming connection")
		} else {
			c.route[connection.RemoteAddr().String()] = connection
			fmt.Printf("%s has joined the server\n", connection.RemoteAddr())
			go c.handle(connection)
		}
	}
}

// Stop the server
func (c *Clan) Stop() { c.listener.Close() }

// Handle incoming connections
func (c *Clan) Handle(conn net.Conn) {
	var dec = gob.NewDecoder(conn)
	var msg Message = Message{}
	for {
		if err := dec.Decode(&msg); err != nil {
			if err == io.EOF {
				delete(c.route, conn.RemoteAddr().String())
				fmt.Printf("%s has left the server: Connection has closed\n", conn.RemoteAddr().String())
				conn.Close()
				return
			}
		}
		fmt.Println("server", msg)
	}
}

func (c *Clan) pass(msg Message) {
	var encoder = gob.NewEncoder(c.route[msg.To])
	if err := encoder.Encode(&msg); err != nil {
		fmt.Println("error: encoding message", err)
	}
}

// Insert data to database
func (c *Clan) Insert(data string) { c.database.Insert(data) }

// Close database
func (c *Clan) Close() { c.database.Close() }

// NewClan start a new server
func NewClan() *Clan {
	var c *Clan = &Clan{
		database: NewWarehouse(),
		route:    make(map[string]net.Conn),
	}
	c.handle = c.Handle
	return c
}
