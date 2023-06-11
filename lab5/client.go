package main

import (
	"encoding/binary"
	"flag"

	"fmt"
	"log"
	"net"
	"time"
)

type Message struct {
	Type           uint8
	Code           uint8
	Checksum       uint16
	Identifier     uint16
	SequenceNumber uint16
	Data           []byte
}

func (m *Message) Marshal() ([]byte, error) {
	length := len(m.Data) + 8
	buf := make([]byte, length)
	buf[0] = m.Type
	buf[1] = m.Code
	buf[2] = 0 // Checksum byte 1
	buf[3] = 0 // Checksum byte 2
	buf[4] = byte(m.Identifier >> 8)
	buf[5] = byte(m.Identifier & 0xff)
	buf[6] = byte(m.SequenceNumber >> 8)
	buf[7] = byte(m.SequenceNumber & 0xff)
	copy(buf[8:], m.Data)
	checksum := Checksum(buf)
	buf[2] = byte(checksum)
	buf[3] = byte(checksum >> 8)
	return buf, nil
}

var ICMP_ECHO_REQUEST uint8 = 8

func Checksum(b []byte) uint16 {
	csumcv := len(b) - 1 // checksum converage
	s := uint32(0)
	// Iterate over 16-bit words
	for i := 0; i < len(b)-1; i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}

	if csumcv&1 == 0 {
		s += uint32(b[csumcv])
	}

	s = s>>16 + s&0xffff
	s = s + s>>16
	return ^uint16(s)
}

func sendOnePing(conn *net.IPConn, ID uint16) error {
	pack := Message{
		Type:           ICMP_ECHO_REQUEST,
		Code:           0,
		Checksum:       0,
		Identifier:     uint16(ID),
		SequenceNumber: 1,
		Data: []byte(fmt.Sprintf("%d", time.Now().UnixNano())),
	}
	bs, err := pack.Marshal()
	if err != nil {
		return err
	}
	_, err = conn.Write(bs)
	return err
}

func receiveOnePing(conn *net.IPConn, ID uint16, timeout uint, dest *net.IPAddr) int64 {
	startedSelect := time.Now()
	err := conn.SetReadDeadline(startedSelect.Add(time.Duration(timeout) * time.Second))
	if err != nil {
		log.Fatalf("Unable to set read deadline: %v\n", err)
	}
	var b [1024]byte
	n, addr, err := conn.ReadFromIP(b[0:])
	if err != nil {
		log.Fatalf("Unable to read from IP: %v\n", err)
	}
	log.Printf("Reply from %v: bytes=%d time=%v\n", addr, n, time.Since(startedSelect))
	receiveTime := time.Now()
	//  Fetch the ICMP header from the IP packet
	m := &Message{
		Type:           uint8(b[0]),
		Code:           uint8(b[1]),
		Checksum:       binary.BigEndian.Uint16(b[2:4]),
		Identifier:     binary.BigEndian.Uint16(b[4:6]),
		SequenceNumber: binary.BigEndian.Uint16(b[6:8]),
		Data:           b[8:n],
	}
	log.Printf("ICMP message: %+v\n", m)
	return time.Duration(receiveTime.Sub(startedSelect)).Milliseconds()
}

func doOnePing(dest *net.IPAddr, timeout int) int64 {
	conn, err := net.DialIP("ip4:icmp", nil, dest)
	if err != nil {
		log.Fatalf("Unable to dial IP: %v\n", err)
	}
	defer conn.Close()
	err = sendOnePing(conn, 0)
	if err != nil {
		log.Fatalf("Unable to send one ping: %v\n", err)
	}
	return receiveOnePing(conn, 0, uint(timeout), dest)
}

func ping(host string, timeout int) {
	// timeout=1 means: If one second goes by without a reply from the server,
	//  the client assumes that either the client's ping or the server's pong is lost
	dest, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		log.Fatalf("Unable to resolve IP address: %v\n", err)
	}
	log.Printf("Pinging %s(%v) with use Golang:\n", host, dest)
	// Send ping requests to a server separated by approximately one second
	delay := doOnePing(dest, timeout)
	log.Printf("Delay: %v\n", delay)
}

var (
	host    string
	timeout int
)

func main() {
	flag.StringVar(&host, "host", "www.baidu.com", "host name")
	flag.IntVar(&timeout, "timeout", 3, "timeout in second")
	flag.Parse()

	ping(host, timeout)
}
