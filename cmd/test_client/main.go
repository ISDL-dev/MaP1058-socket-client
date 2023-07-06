// a client to check the connection with a server
package main

import (
	"fmt"
	"log"

	"github.com/Be3751/socket-capture-signals/internal/pkg/socket"
)

func main() {
	conf := socket.Config{
		ServerIP:   "localhost",
		ServerPort: "3000",
		ClientIP:   "localhost",
		ClientPort: 1000,
	}
	conn, err := socket.Connect(conf)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write([]byte("test"))
	if err != nil {
		log.Fatal(err)
	}
	receiveBuf := make([]byte, 1024)
	_, err = conn.Read(receiveBuf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(receiveBuf))
}
