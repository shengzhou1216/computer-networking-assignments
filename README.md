# Computer Netowrking: A TopDown Approach - Assignments(Golang)

http://gaia.cs.umass.edu/kurose_ross/programming.php

- [x] Lab1: A simple client/server simple socket program
- [x] Lab2: Web server
- [x] Lab3: UDP	Pinger
- [x] Lab4: SMTP
- [x] Lab5: ICMP pinger
- [x] Lab6: HTTP Web Proxy Server

## Lab1: A simple client/server simple socket program
In this assignment, you’ll write a client that will use sockets to communicate with a server that you will also write.  Here’s what your client and server should do:

Your client should first accept an integer between 1 and 100 from the keyboard,  open a TCP socket to your server and send a message containing (i) a string containing your name (e.g., “Client of John Q. Smith”) and (ii) the entered integer value and then wait for a sever reply.

## Lab2: Web server
In this lab, you will learn the basics of socket programming for TCP connections in Python: how to create 
a socket, bind it to a specific address and port, as well as send and receive a HTTP packet. You will also 
learn some basics of HTTP header format.
You will develop a web server that handles one HTTP request at a time. Your web server should accept 
and parse the HTTP request, get the requested file from the server’s file system, create an HTTP response 
message consisting of the requested file preceded by header lines, and then send the response directly to 
the client. If the requested file is not present in the server, the server should send an HTTP “404 Not 
Found” message back to the client

## Lab3: UDP Pinger
In this lab, you will learn the basics of socket programming for UDP in Python. You will learn how to send and receive datagram packets using UDP sockets and also, how to set a proper socket timeout. 
Throughout the lab, you will gain familiarity with a Ping application and its usefulness in computing  statistics such as packet loss rate

## Lab4: SMTP
By the end of this lab, you will have acquired a better understanding of SMTP protocol. You will also 
gain experience in implementing a standard protocol using Python.
Your task is to develop a simple mail client that sends email to any recipient. Your client will need to 
connect to a mail server, dialogue with the mail server using the SMTP protocol, and send an email 
message to the mail server. Python provides a module, called smtplib, which has built in methods to send 
mail using SMTP protocol. However, we will not be using this module in this lab, because it hide the 
details of SMTP and socket programming

## Lab5: ICMP pinger
In this lab, you will gain a better understanding of Internet Control Message Protocol (ICMP). You will  learn to implement a Ping application using ICMP request and reply messages.

## Lab6: HTTP Web Proxy Server
In this lab, you will learn how web proxy servers work and one of their basic functionalities –caching. 
Your task is to develop a small web proxy server which is able to cache web pages. It is a very simple proxy server which only understands simple GET-requests, but is able to handle all kinds of objects -
not just HTML pages, but also images.