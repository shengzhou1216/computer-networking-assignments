package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

func main() {
	l, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 5555,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on %s\n", l.LocalAddr().String())
	defer l.Close()
	var buf [1024]byte
	for {
		n, remoteAddr, err := l.ReadFromUDP(buf[0:])
		if err != nil {
			log.Println(err)
		}
		random := rand.Intn(10)
		if random < 4 {
			continue
		}
		data := string(buf[:n])
		seq := strings.Split(data, " ")[1]

		log.Printf("%s: %s\n",remoteAddr, strings.TrimRight(data, "\n"))
		time.Sleep(time.Duration(random) * time.Millisecond)
		l.WriteToUDP([]byte(fmt.Sprintf("Pong %s\n", seq)), remoteAddr)
	}

}
