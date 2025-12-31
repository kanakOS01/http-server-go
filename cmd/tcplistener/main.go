package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		defer f.Close()
		defer close(out)

		buf := make([]byte, 8)
		var line strings.Builder

		for {
			n, err := f.Read(buf)

			if n > 0 {
				parts := strings.Split(string(buf[:n]), "\n")
				lenParts := len(parts)

				for i, part := range parts {
					line.WriteString(part)
					if i != lenParts-1 {
						out <- line.String()
						line.Reset()
					}
				}
			}

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("read error:", err)
				return
			}
		}

		if line.Len() > 0 {
			out <- line.String()
		}
	}()

	return out
}

func main() {
	// set up a listener for tcp conn at port 42069
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		// anything that hits the part comes here (stream of data)
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())

		linesChan := getLinesChannel(conn)
		for line := range linesChan {
			fmt.Println(line)
		}

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}
