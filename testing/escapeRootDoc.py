
from socket import socket

# Create connection to the server

s = socket()

s.connect(("localhost", 8080));

# Compose the message/HTTP request we want to send to the server

# test if url escape root doc, this should be 200 ok
msgPart1 = b"GET /subdir1/../ HTTP/1.1\r\nHost: Ha11\r\n\r\n"

# Send out the request

s.sendall(msgPart1)

# test if url escape root doc, this should be 404
msgPart2 = b"GET /subdir100/../ HTTP/1.1\r\nHost: Ha11\r\n\r\n"
s.sendall(msgPart2)

# Listen for response and print it out

print (s.recv(4096))

s.close()