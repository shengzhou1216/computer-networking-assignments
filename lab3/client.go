package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	remoteAddr :=  &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 5555,
	}
	conn,err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	
	for i := 1; i <= 10; i++ {
		now := time.Now().UnixMilli()
		ping := fmt.Sprintf("Ping %d %d\n", i,now)
		_,err := conn.Write([]byte(ping))
		if err != nil {
			log.Fatal(err)
		}
		log.Println(ping)
		// set timeout
		err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			log.Fatal(err)
		}
		data,err := bufio.NewReader(conn).ReadString('\n')
		data = strings.TrimRight(data, "\n")
		if err != nil {
			log.Println("Request timed out")
			continue
		}
		// calculate RTT in ms
		rtt := time.Now().UnixMilli() - now
		log.Println(fmt.Sprintf("Data: %s RTT: %d ms",data,rtt))
	}
		
}