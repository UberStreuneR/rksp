package main

import (
	// "practice04/controllers"
	"log"
	"practice04/initializers"
	// "github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables\n", err.Error())
	}
	initializers.ConnectDB(&config)
}

func main() {
	// s := server.NewRsocketServer("rofl", *initializers.DB)
	// go s.Serve()
	// server.Rsocket_Client()
}
