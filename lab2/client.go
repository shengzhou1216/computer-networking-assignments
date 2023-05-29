package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"fmt"
)

func main() {
	var host string
	var port int
	var file string
	flag.StringVar(&host,"host", "localhost", "host")
	flag.IntVar(&port, "port", 8080, "port")
	flag.StringVar(&file, "file", "/", "file")
	flag.Parse()

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// send request
	conn.Write([]byte(fmt.Sprintf("GET %s HTTP/1.1\r\n", file)))
	conn.Write([]byte(fmt.Sprintf("Host: %s:%d\r\n", host,port)))
	conn.Write([]byte("\r\n"))
	// read response
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		log.Println(scanner.Text())
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
}
