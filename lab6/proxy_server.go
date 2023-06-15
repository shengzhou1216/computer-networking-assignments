package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var (
	proxyAddr = flag.String("addr", ":8000", "proxy server address")
)

func main() {
	flag.Parse()
	l,err := net.Listen("tcp", *proxyAddr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	log.Println("proxy server listen on", *proxyAddr)
	for {
		conn,err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer func ()  {
		if err := recover(); err != nil {
			log.Println("handleConn panic:", err)
		}
	}()
	log.Printf("Received a connection from %s\n", conn.RemoteAddr())
	defer conn.Close()
	reader := bufio.NewReader(conn)
	// request line
	line,err := reader.ReadString('\n')
	if err != nil {
		log.Println("read error: ", err)
		return
	}
	parts := strings.Split(line," ")
	if len(parts) != 3 {
		log.Println("invalid request line:", line)
		return
	}
	method,url,version := parts[0],parts[1],parts[2]
	log.Printf("method: %s, url: %s, version: %s",method,url,version)
	// header lines
	headers := make(map[string]string)
	for {
		line,err := reader.ReadString('\n')
		if err != nil {
			log.Println("read error:", err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.Split(line,":")
		if len(parts) == 2 {
			headers[parts[0]] = parts[1]
		}
	}
	// log.Println("headers:", headers)
	//  Extract the filename from the given message
	filename := strings.SplitN(url,"/",2)[1]
	log.Println("filename:", filename)
	// check wether the file exist in the cache
	if _,err := os.Stat(filename); err == nil || os.IsExist(err) {
		log.Println("File found in cache")
		// ProxySever finds a cache hit and generates a response message
		conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
		conn.Write([]byte("Content-Type: text/html\r\n"))
		// Blank line
		conn.Write([]byte("\r\n"))
		// Read the file from disk into the response message
		f,err := os.OpenFile(filename,os.O_RDONLY,0666)
		if err != nil {
			log.Println("open file error:", err)
			return
		}
		defer f.Close()
		bs,err := io.ReadAll(f)
		if err != nil {
			log.Println("read file error:", err)
			return
		}
		// Send file
		conn.Write(bs)
		conn.Write([]byte("\r\n"))
		log.Println("Read from cache")
	}else { // Error handling for file not found in cache
		log.Println("File not found in cache")
		// Create a socket on the proxyserver
		hostn := strings.Replace(filename,"www.","",1)
		log.Println("hostn:", hostn)
		// connect to the socket to port 80
		c,err := net.Dial("tcp",hostn+":80")
		if err != nil {
			log.Println("dial error:", err)
			return
		}
		defer c.Close()
		// Create a temporary file on the socket and ask port 80 for the file requested by the client
		c.Write([]byte("GET " + "http://" + filename + " HTTP/1.1\r\n"))
		c.Write([]byte("Host: " + hostn+":80" + "\r\n"))
		c.Write([]byte("Connection: close\r\n"))
		c.Write([]byte("\r\n"))
		log.Println("Write to server")
		// Read the response into buffer
		// err = c.SetWriteDeadline(time.Now().Add(5 * time.Second))
		// if err != nil {
		// 	log.Println("set write deadline error:", err)
		// }
		bs,err := io.ReadAll(c)
		if err != nil {
			log.Println("read error:", err)
			return
		}
		// Create a new file in the cache for the requested file.
		tmpFile,err  := os.OpenFile("./" + filename,os.O_WRONLY|os.O_CREATE,0666)
		if err != nil {
			log.Println("open error:", err)
			return
		}
		defer tmpFile.Close()
		// Write the response body to this file
		s := string(bs)
		body := strings.Split(s,"\r\n\r\n")[1]
		tmpFile.Write([]byte(body))
		// Also send the response in the buffer to client socket and the corresponding file in the cache
		c.Write(bs)
		log.Println("Read from server")
	}
}