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
	var cNumber int
	var cName = "Sheng"

	for {
		fmt.Print("Please input the number: ")

		_, err := fmt.Scan(&cNumber)

		if err != nil {
			fmt.Println("Invalid input. Please enter a valid number.")
			continue
		}
		break
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	if _, err := conn.Write([]byte(fmt.Sprintf("Name:%s\nNumber:%d\n\r\n", cName, cNumber))); err != nil {
		log.Fatal(err)
	}
	conn.CloseWrite()
	reader := bufio.NewReader(conn)
	var sNumber int
	var sName string
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(data, "\r\n") {
			break
		}
		data = strings.TrimSuffix(data, "\n")
		log.Println(fmt.Sprintf("Client received [%s] from server", data))
		if strings.Contains(data, "Name:") {
			sName = strings.ReplaceAll(data, "Name:", "")
		} else if strings.Contains(data, "Number:") {
			sNumber, _ = strconv.Atoi(strings.ReplaceAll(data, "Number:", ""))
		}
	}
	log.Println(fmt.Sprintf("Client: %s, Number: %d; Server: %s, Number: %d. Sum: %d", cName, cNumber, sName, sNumber, cNumber+sNumber))
}
