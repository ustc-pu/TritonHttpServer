# test bad concurrent
from socket import socket
import time
# Create connection to the server

s = socket()

s.connect(("localhost", 8080));

# Compose the message/HTTP request we want to send to the server

msgPart1 = b"GET /subdir1 HTTP/1.1\r\nHost: Ha11\r"

# Send out the request

s.sendall(msgPart1)


msgPart2 = b"\n\r\n"
time.sleep(5.1)
s.sendall(msgPart2)
# Listen for response and print it out

print (s.recv(4096))

s.close()