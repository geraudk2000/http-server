package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

var directory string

func respond(conn net.Conn, message string) {
	conn.Write([]byte(message))
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	requestLine := strings.Split(string(buf[:n]), "\r\n")[0]
	parts := strings.Split(requestLine, " ")
	if len(parts) < 2 {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}
	method := parts[0]
	path := parts[1]

	switch {
	case method == "GET" && path == "/":
		respond(conn, "HTTP/1.1 200 OK\r\n\r\n")
	case method == "GET" && strings.HasPrefix(path, "/echo/"):
		body := strings.TrimPrefix(path, "/echo/")
		contentLength := len(body)
		respond(conn, fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentLength, body))
	case method == "GET" && path == "/user-agent":
		// user-agent response
		userAgent := ""
		lines := strings.Split(string(buf[:n]), "\r\n")

		for _, line := range lines {
			if strings.HasPrefix(line, "User-Agent:") {
				userAgent = strings.TrimSpace(strings.TrimPrefix(line, "User-Agent:"))
				break
			}
		}
		contentLength := len(userAgent)
		respond(conn, fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentLength, userAgent))

	case strings.HasPrefix(path, "/files/"):
		filename := strings.TrimPrefix(path, "/files")
		filePath := filepath.Join(directory, filename)

		if method == "GET" {
			file, err := os.Open(filePath)
			if err != nil {
				respond(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
			} else {
				defer file.Close()
				content, err := io.ReadAll(file)
				if err != nil {
					respond(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\n")
				} else {
					contentLength := len(content)
					respond(conn, fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", contentLength, content))
				}
			}
		} else if method == "POST" {
			// handle file upload
			headersBody := strings.SplitN(string(buf[:n]), "\r\n\r\n", 2)
			headers := strings.Split(headersBody[0], "\r\n")
			body := headersBody[1]
			contentLength := 0

			for _, line := range headers {
				if strings.HasPrefix(line, "Content-Length:") {
					lengthStr := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
					contentLength, err = strconv.Atoi(lengthStr)
					if err != nil {
						respond(conn, "HTTP/1.1 400 Bad Request\r\n\r\n")
						return
					}
					break
				}
			}
			//fmt.Println(contentLength, body)

			for len(body) < contentLength {
				moreBuf := make([]byte, contentLength-len(body))
				m, err := conn.Read(moreBuf)
				if err != nil {
					respond(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\n")
					return
				}
				body += string(moreBuf[:m])
			}
			err := os.WriteFile(filePath, []byte(body), 0644)
			if err != nil {
				respond(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\n")
			} else {
				respond(conn, "HTTP/1.1 201 Created\r\n\r\n")
			}
		}
	default:
		respond(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
	}
}

func main() {
	//response = ""

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	dirFlag := flag.String("directory", ".", "directory to serve files from")
	flag.Parse()
	directory = *dirFlag

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
