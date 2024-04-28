package main

import (
	"log"

	"github.com/tmavrin/tp-link/client"
)

func main() {
	c, err := client.Authenticate("http://192.168.1.1", "username", "password")
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()
}
