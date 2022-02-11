package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	server := NewChatServer()
	handler := &ConnectionHandler{
		Server: server,
	}
	r := gin.Default()
	r.GET("/", handler.HandleConnection)
	log.Fatal(r.Run(":9000"))
}
