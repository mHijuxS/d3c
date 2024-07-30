package main

import (
	"commons/commons"
	"encoding/gob"
	"log"
	"net"
)

var (
	fieldAgents = []commons.Message{}
)

func main() {
	log.Println("Started Execution")
	startListener("9090")
}

func agentIsRegistered(agentId string) (registered bool) {
	registered = false

	for _, agent := range fieldAgents {
		if agent.AgentId == agentId {
			registered = true
			break
		}
	}
	return registered
}

func messageContainsResponse(message commons.Message) (containsResponse bool) {
	containsResponse = false

	for _, command := range message.Commands {
		if len(command.Response) > 0 {
			containsResponse = true
			break
		}
	}
	return containsResponse
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

				// Verify if the agent is already registered
				if agentIsRegistered(message.AgentId) {
					log.Println("Message from Agent: ", message.AgentId)
					if messageContainsResponse(*message) {
						// Print the response
						for _, command := range message.Commands {
							log.Println("Command: ", message.Commands)
							log.Println("Response: ", command.Response)
						}
					}
				} else {
					log.Println("Registering Agent: ", message.AgentId)
					fieldAgents = append(fieldAgents, *message)
				}
				gob.NewEncoder(channel).Encode(message)
			}

		}

	}
}
