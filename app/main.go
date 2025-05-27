package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit
var response string

func main() {
	response := ""

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	//
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
	}
	//fmt.Println(string(buf[:n]))
	// Parse request line

	requestLine := strings.Split(string(buf[:n]), "\r\n")[0]
	parts := strings.Split(requestLine, " ")

	if len(parts) >= 2 && parts[1] == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
	} else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	//response := "HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n"
	conn.Write([]byte(response))
	conn.Close()
}
