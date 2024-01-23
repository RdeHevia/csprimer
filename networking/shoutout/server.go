package main

import (
	"log"
	"net"
	"strings"
)

func main() {
	connection, err := net.ListenPacket("udp4", "localhost:1234")
	if err != nil {
		log.Fatalln(err)
	}

	defer connection.Close()

	buffer := make([]byte, 1024)

	log.Println("Listening on port 1234")

	for {
		byteCount, callingAddress, err := connection.ReadFrom(buffer)
		if err != nil {
			log.Fatalln(err)
		}

		if byteCount > 0 {
			log.Println("byteCount", byteCount)
			log.Println("callingAddress", callingAddress)
			messageReceived := strings.TrimSpace(string(buffer[0 : byteCount-1]))
			log.Println("Message received:", messageReceived)

			if messageReceived == "STOP" {
				log.Println("Exiting UDP server!")
				return
			}
			response := "Hello! I have received this message: " + messageReceived
			responseBytes := []byte(response)
			_, err = connection.WriteTo(responseBytes, callingAddress)
		}
	}
}
