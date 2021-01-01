package tritonhttp

import (
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"time"
)
const readTimeout = 5 * time.Second
/* 
For a connection, keep handling requests until 
	1. a timeout occurs or
	2. client closes connection or
	3. client sends a bad request
*/
// Start a loop for reading requests continuously
// Set a timeout for read operation
// Read from the connection socket into a buffer
// Validate the request lines that were read
// Handle any complete requests
// Update any ongoing requests
// If reusing read buffer, truncate it before next read

func (hs *HttpServer) handleConnection(conn net.Conn) {
	log.Println("Start handling new connection.")
	defer conn.Close()
	defer log.Println("Closed connection.")
	connectionSetClose := false
	for !connectionSetClose {
		setTimeOutErr := conn.SetReadDeadline(time.Now().Add(readTimeout)) // timeout
		if setTimeOutErr != nil {
			log.Println("setReadDeadline failed:", setTimeOutErr)
		} else {
			log.Println("Set up a read timeout for 5 s")
		}
		header := ""
		size := 0
		readErr := errors.New("")
		for {
			buf := make([]byte, 1024)
			size, readErr = conn.Read(buf)
			if readErr != nil {
				if netErr, ok := readErr.(net.Error); ok && netErr.Timeout() {
					log.Println("time out error", readErr)
					if len(header) > 0 {
						httpRequestHeader, _ := newHttpRequestHeader(header)
						hs.handleBadRequest(httpRequestHeader, conn)
					}
				}
				return
				//else if readErr == io.EOF {
				//	log.Println("Read to the end of the client request: ", readErr)
				//	if !connectionSetClose {
				//		log.Println("client did not set connection close, wait for", readTimeout)
				//		time.Sleep(readTimeout)
				//	}
				//	break
				//} else {
				//	log.Println("Read failed: ", readErr)
				//	break
				//}
			}
			data := buf[:size]
			header += string(data)
			//check if header contains a request e.g. "\r\n\r\n", that's the end of a full request
			if idx := strings.Index(header, "\r\n\r\n"); idx > 0 {
				httpRequestHeader, buildHeaderErr := newHttpRequestHeader(header[:idx+4])
				validReq := true
				if httpRequestHeader.requestHeaderMap["Connect"] == "close" {
					log.Println("client set connection close")
					connectionSetClose = true
				}
				if buildHeaderErr != nil {
					log.Println("Found error when building new requestHeader:", buildHeaderErr)
					if buildHeaderErr.Error() == "malformed request header" {
						validReq = hs.validateRequestHeader(httpRequestHeader, conn, false)
					}
				} else {
					validReq = hs.validateRequestHeader(httpRequestHeader, conn, true)
				}
				if !validReq {
					connectionSetClose = true
					break
				}
				header = header[idx+4:] //update header
				log.Print("get next request...")
			}
		}
	}
}

//
func (hs *HttpServer) validateRequestHeader(requestHeader *HttpRequestHeader, conn net.Conn, isValid bool) (validReq bool) {

	log.Println("now validate httpHeader...")
	var iniLine error
	if !isValid {
		log.Println("invalid http request")
		iniLine = errors.New("invalid http request")
		hs.handleBadRequest(requestHeader, conn)
		return false
	}

	if requestHeader.method != "GET" {
		log.Println("unsupported method " + requestHeader.method + " found")
		iniLine = errors.New("invalid methods")
		hs.handleBadRequest(requestHeader, conn)
		return false
	}

	if requestHeader.version != "HTTP/1.1" {
		log.Println("unsupported http version " + requestHeader.version + " found")
		hs.handleBadRequest(requestHeader, conn)
		iniLine = errors.New("invalid protocol")
		return false
	}

	if requestHeader.url[0] != '/' || !strings.Contains(requestHeader.url, "/") {
		log.Println("invalid http url " + requestHeader.url + " found")
		hs.handleBadRequest(requestHeader, conn)
		iniLine = errors.New("invalid url")
		return false
	}

	if _, ok := requestHeader.requestHeaderMap["Host"]; !ok {
		log.Println("invalid request header, does not contain Host")
		hs.handleBadRequest(requestHeader, conn)
		iniLine = errors.New("invalid request header")
		return false
	}

	if iniLine == nil {
		//check if path exist
		_, err := os.Stat(hs.DocRoot + requestHeader.url)
		if err != nil {
			log.Println("file not found on server...")
			hs.handleFileNotFoundRequest(requestHeader, conn)
			log.Println(err)
			return true
		} else {
			log.Println("file found on server...")
			hs.handleResponse(requestHeader, conn)
			return true
		}
	}
	return true
}
