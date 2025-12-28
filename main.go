package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

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
