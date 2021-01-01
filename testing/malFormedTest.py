from socket import socket
# Create connection to the server

s = socket()

s.connect(("localhost", 8080));

# Compose the message/HTTP request we want to send to the server
# test malformed req, ":" missing 400 bad req
msgPart1 = b"GET /index.html HTTP/1.1\r\nHost Ha11\r\n\r\n"

# Send out the request
s.sendall(msgPart1)

# test malformed req, kv pair invalid, 400 bad req
msgPart2 = b"GET /index.html HTTP/1.1\r\nHost:\r\n\r\n"
s.sendall(msgPart2)

# Listen for response and print it out

print (s.recv(4096))

s.close()