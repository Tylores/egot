package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	log.Println("Starting client access test")

	server := exec.Command("./bin/core")
	server.Stdout = os.Stdout
	server.Stderr = os.Stderr

	err := server.Start()
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
	time.Sleep(100 * time.Millisecond)

	client := exec.Command("./bin/crawler")
	client.Stdout = os.Stdout
	client.Stderr = os.Stderr

	err = client.Run()
	if err != nil {
		log.Fatalf("Failed to run client: %v", err)
	}

	err = server.Process.Kill()
	if err != nil {
		log.Fatalf("Failed to stop server: %v", err)
	}

	log.Println("Finised client access test")
}
