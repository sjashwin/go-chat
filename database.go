package chat

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// Database that will be used to store the messagesa
type Database interface {
	Create(filepath string)
	Insert(data string) error
	Close()
}

// Warehouse is a type of database that stores the messages
// in a simple log file
type Warehouse struct {
	file *os.File
}

// Create a new database
func (w *Warehouse) Create(filepath string) {
	os, err := os.Create(filepath)
	if err != nil {
		panic("error: could not create a new file")
	}
	w.file = os
}

// Insert data to the database
func (w Warehouse) Insert(data string) error {
	var current time.Time = time.Now()
	var format string = fmt.Sprintf("%v %s\n", current, data)
	var length, err = w.file.WriteString(format)
	if len(data) <= length {
		return errors.New("error: did not write data correctly")
	}
	return err
}

// Close the database
func (w *Warehouse) Close() {
	w.file.Close()
}

// NewWarehouse is a new database
func NewWarehouse() *Warehouse {
	var warehouse *Warehouse = &Warehouse{}
	return warehouse
}
