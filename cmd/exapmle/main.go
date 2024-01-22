package main

import (
	client "github.com/ISDL-dev/MaP1058-socket-client"
	"time"
)

const (
	outputDir = "output"
)

func main() {
	conf := client.Config{
		ServerIP: "192.168.10.128",
		SaveDir:  outputDir,
	}

	c, err := client.NewClient(conf)
	if err != nil {
		panic(err)
	}
	err = c.Start(time.Second * 3)
	if err != nil {
		panic(err)
	}
}
