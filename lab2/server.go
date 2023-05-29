package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	log.Println(fmt.Sprintf("Server listen on: %s", l.Addr().String()))
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	// request line
	parts := strings.Split(line, " ")
	method := parts[0]
	url := parts[1]
	protocol := parts[2]
	log.Println("Method:", method)
	log.Println("Url:", url)
	log.Println("Protocol:", protocol)
	// header lines
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.Split(line, ": ")
		if len(parts) == 2 {
			headers[parts[0]] = parts[1]
		}
	}
	log.Println("Headers:", headers)
	// check if file exists
	wd, _ := os.Getwd()
	filePath := path.Join(wd, url)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n"))
		conn.Write([]byte("Content-Type: text/html\r\n"))
		conn.Write([]byte("\r\n"))
		conn.Write([]byte("404 Not Found"))
		return
	}
	// read file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// write response
	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Write([]byte("Content-Type: text/html\r\n"))
	conn.Write([]byte("\r\n"))
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		conn.Write([]byte(line))
	}
}
