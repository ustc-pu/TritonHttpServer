# test pipelined requests, where second is 404 file not found, and the first and third are valid.

from socket import socket

# Create connection to the server

s = socket()

s.connect(("localhost", 8080));

# Compose the message/HTTP request we want to send to the server
msgPart1 = b"GET / HTTP/1.1\r\nHost: Ha11\r\n\r\n"

# Send out the request
s.sendall(msgPart1)

# Listen for response and print it out
print (s.recv(4096))

msgPart2 = b"GET /index100.html HTTP/1.1\r\nHost: Ha22\r\n\r\n"
s.sendall(msgPart2)
print (s.recv(4096))

msgPart2 = b"GET /index.html HTTP/1.1\r\nHost: Ha22\r\n\r\n"
s.sendall(msgPart2)
print (s.recv(4096))

s.close()