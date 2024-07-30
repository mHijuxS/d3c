package main

import (
	"log"
	"os"
	"time"
	"crypto/md5"
	"encoding/hex"
)

func main() {
	log.Println("Started Execution")
	log.Println("Agent ID: %s", makeId())
}

func makeId() string{
	myHostname,_ := os.Hostname()
	myTime := time.Now().String()

	hasher := md5.New()
	hasher.Write([]byte(myHostname + myTime))

	return hex.EncodeToString(hasher.Sum(nil))
}