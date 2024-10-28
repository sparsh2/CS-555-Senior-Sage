package main

// import gin

import (
	"fmt"

	"authz/config"
	"authz/web"
)

func main() {
	// Load the configs
	config.LoadConfig()
	// Start the server
	r := web.GetRouter()
	err := r.Run()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
