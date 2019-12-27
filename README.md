# A Chat Application in golang

A single server that routes messages to multiple sources

## Getting Started Server

```go
package main

import(
    "fmt"
    "go-chat"
)

func main(){
    serve := chat.NewClan()
    fmt.Println("Starting server...")
}
```

## Getting Started Client

```go
package main

import "go-chat"

func main(){
    chat.NewClient("name")
}
```