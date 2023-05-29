package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":8080")
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		l.Close()
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println(fmt.Sprintf("Server is listening on port %s", addr))
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	log.Println(fmt.Sprintf("Server accepted a connection from %s", remoteAddr))
	var sName = "ShengServer"
	var sNumber int
	var cName string
	var cNumber int
	reader := bufio.NewReader(conn)
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(data,"\r\n") {
			break
		}
		data = strings.TrimSuffix(data, "\n")
		log.Println(fmt.Sprintf("Server received [%s] from %s", data, remoteAddr))
		if strings.Contains(data, "Name:") {
			cName = strings.ReplaceAll(data, "Name:", "")
		} else if strings.Contains(data, "Number:") {
			cNumber, _ = strconv.Atoi(strings.ReplaceAll(data, "Number:", ""))
			if cNumber < 0 || cNumber > 100 {
				log.Println(fmt.Sprintf("Client %s sent an invalid number %d. Close the connection.", cName, cNumber))
				conn.Close()
				return
			}
			sNumber = cNumber * 2
		}
	}
	log.Println(fmt.Sprintf("Client: %s, Number: %d; Server: %s, Number: %d", cName, cNumber, sName, sNumber))
	conn.Write([]byte(fmt.Sprintf("Name:%s\nNumber:%d\n\r\n", sName, sNumber)))
	conn.CloseWrite()
}
