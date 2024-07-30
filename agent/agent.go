package main

import (
	"commons/commons"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"log"
	"net"
	"os"
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

		gob.NewEncoder(channel).Encode(message)
		gob.NewDecoder(channel).Decode(&message)

		time.Sleep(time.Duration(waitTime) * time.Second)
	}

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
