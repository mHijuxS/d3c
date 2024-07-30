package main

import ("log"
		"net")

func main(){
	log.Println("Started Execution")
	startListener("9090")
}

func startListener(port string){
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("Error starting listener: %v", err)
	}else{
		log.Println("Listening on port %s", port)
		channel, err := listener.Accept()
		defer channel.Close()

		if err != nil {
			log.Fatalf("Error accepting connection: %v", err)
		}
		
		log.Println("Accepted connection from %s", channel.RemoteAddr())

	}
}