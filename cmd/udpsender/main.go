package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// resolve UDP address
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal("Error resolving UPD address: ", err)
	}

	// prepare a UPD connection (not actually created)
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("Error preparing UDP connection: ", err)
	}
	defer conn.Close()

	// read from stdin
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading from stdin: ", err)
			continue
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Println("Error writing to UDP: ", err)
			continue
		}
	}

}