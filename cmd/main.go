package main

import (
	client "github.com/Be3751/MaP1058-socket-client"
	"time"
)

const (
	outputDir = "output"
)

func main() {
	conf := client.Config{
		ServerIP: "192.168.10.101",
		SaveDir:  outputDir,
	}

	c, err := client.NewClient(conf)
	if err != nil {
		panic(err)
	}
	err = c.Start(time.Second * 10)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 5)
	err = c.Stop()
	if err != nil {
		panic(err)
	}
}
