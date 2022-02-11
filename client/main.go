package main

import (
	"context"
)

func main() {
	cli := &Client{
		ServerURL: "ws://localhost:9000",
		Ctx:   context.Background(),
	}
	cli.Run()
}
