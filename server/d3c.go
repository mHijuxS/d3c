package main

import (
	"commons/commons"
	"encoding/gob"
	"log"
	"net"
)

func main() {
	log.Println("Started Execution")
	startListener("9090")
}

func startListener(port string) {
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("Error starting listener: %v", err)
	} else {
		log.Printf("Listening on port %s", port)
		for {
			channel, err := listener.Accept()
			defer channel.Close()

			if err != nil {
				log.Fatalf("Error accepting connection: %v", err)
			} else {
				message := &commons.Message{}

				gob.NewDecoder(channel).Decode(message)

				log.Println("Agent Message: ", message.AgentId)
				log.Printf("New connection from %s", channel.RemoteAddr())

				//
				//
				//

				gob.NewEncoder(channel).Encode(message)
			}

		}

	}
}
