package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	from := "a@b.com"
	to := "x@y.com"
	authCode := "<auth code for <from>>"
	subject := "Testing"

	// Choose a mail server(e.g. Google mail server) and call it mailserver
	mailserver := "smtp.sina.com:25"
	raddr, err := net.ResolveTCPAddr("tcp", mailserver)
	if err != nil {
		log.Fatal(err)
	}
	// Create socket called clientSocket and establish a TCP connection with mailserver
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	recv := getReply(conn)
	assertReplyCode(recv, "220")

	// Send HELO command and print server response.
	send(conn, "HELO Alice\r\n")
	recv = getReply(conn)
	assertReplyCode(recv, "250")

	// AUTH LOGIN
	send(conn, "AUTH LOGIN\r\n")
	recv = getReply(conn)
	assertReplyCode(recv, "334")
	
	// Send username
	send(conn, fmt.Sprintf("%s\r\n", base64.StdEncoding.EncodeToString([]byte(from))))
	recv = getReply(conn)
	assertReplyCode(recv, "334")

	// Send password/auth code
	send(conn, fmt.Sprintf("%s\r\n", base64.StdEncoding.EncodeToString([]byte(authCode))))
	recv = getReply(conn)
	assertReplyCode(recv, "235")


	// Send MAIL FROM command amd print server response.
	fromCommand := fmt.Sprintf("MAIL FROM:<%s>\r\n",from)
	send(conn, fromCommand)
	recv = getReply(conn)
	assertReplyCode(recv, "250")

	// Send RCPT TO command and print server response.
	rcptCommand := fmt.Sprintf("RCPT TO:<%s>\r\n",to)
	send(conn, rcptCommand)
	recv = getReply(conn)
	assertReplyCode(recv, "250")
	
	// Send DATA command and print se%rver response.
	dataCommand := fmt.Sprintf("DATA \r\n",)
	send(conn, dataCommand)
	recv = getReply(conn)
	assertReplyCode(recv, "354")

	// Send message data.
	send(conn, "Date: " + time.Now().Format(time.RFC822) + "\r\n")
	send(conn, "From: " + from + "\r\n")
	send(conn, "To: " + to + "\r\n")
	send(conn, "Subject: " + subject + "\r\n")
	send(conn, "\r\n") // empty line to separate headers from body, see RFC5322
	send(conn, "I love computer networks!\r\n")
	send(conn, "I love computer networks!\r\n")
	// Message ends with a single period.
	send(conn, "\r\n.\r\n")

	recv = getReply(conn)
	assertReplyCode(recv, "250")

	// Send QUIT command and get server response.
	send(conn, "QUIT\r\n")
	recv = getReply(conn)
	assertReplyCode(recv, "221")
}

func send(conn *net.TCPConn, command string) {
	conn.Write([]byte(command))
	log.Printf("Send: %s\n", command)
}

func getReply(conn *net.TCPConn) string {
	recv, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Reply: %s\n", recv)
	return recv
}

func assertReplyCode(reply, code string) {
	if reply[0:3] != code {
		log.Printf("%s reply not received from server.\n",code)
	}
}