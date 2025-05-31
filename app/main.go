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

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var response string
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
	}
	requestLine := strings.Split(string(buf[:n]), "\r\n")[0]
	parts := strings.Split(requestLine, " ")
	if len(parts) >= 2 {
		path := parts[1]

		if path == "/" {
			response = "HTTP/1.1 200 OK\r\n\r\n"
		} else if strings.HasPrefix(path, "/echo/") {
			body := strings.TrimPrefix(path, "/echo/")
			contentLength := len(body)

			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentLength, body)
		} else if strings.HasPrefix(path, "/user-agent") {
			userAgent := ""
			lines := strings.Split(string(buf[:n]), "\r\n")
			//fmt.Println(lines)

			for _, line := range lines {
				if strings.HasPrefix(line, "User-Agent:") {
					userAgent = strings.TrimSpace(strings.TrimPrefix(line, "User-Agent:"))
					break
				}
			}
			contentLength := len(userAgent)
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentLength, userAgent)

		} else {
			response = "HTTP/1.1 404 Not Found\r\n\r\n"
		}
	} else {
		response = "HTTP/1.1 400 Bad Request\r\n\r\n"
	}

	//response := "HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n"
	conn.Write([]byte(response))

}

func main() {
	//response = ""

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	// handle multiple connection

	for {

		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}

}
