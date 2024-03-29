// a client to check the connection with a server
package main

import (
	"fmt"
	"log"

	"github.com/ISDL-dev/MaP1058-socket-client/internal/socket"
)

func main() {
	conf := socket.SocketConfig{
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
