package main

import (
	"fmt"

	"github.com/raynine/parking-go/parking"
)

func main() {
	server := parking.NewServer()
	server.Init()

	fmt.Printf("Server starting on port: 8080")

}
