package tritonhttp

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func (hs *HttpServer) handleBadRequest(requestHeader *HttpRequestHeader, conn net.Conn) {
	//panic("todo - handleBadRequest")
	var header string = ""
	header = buildRespContent("HTTP/1.1", "400", "Bad Request")
	header += hs.buildResponseMap("400", requestHeader.url, false)
	//responseHeader := newHttpResponseHeader(header)
	//hs.sendResponse(content, *responseHeader, conn)
	conn.Write([]byte(header))
}

func (hs *HttpServer) handleFileNotFoundRequest(requestHeader *HttpRequestHeader, conn net.Conn) {
	//panic("todo - handleFileNotFoundRequest")
	header := ""
	header = buildRespContent(requestHeader.version, "404", "Not Found")
	header += hs.buildResponseMap("404", requestHeader.url, false)
	//responseHeader := newHttpResponseHeader(content)
	//hs.sendResponse(content, *responseHeader, conn)
	conn.Write([]byte(header))
}

func (hs *HttpServer) handleResponse(requestHeader *HttpRequestHeader, conn net.Conn)  {
	//panic("todo - handleResponse")
	header := ""
	header = buildRespContent(requestHeader.version, "200", "OK")

	header += hs.buildResponseMap("200", requestHeader.url, true)
	responseHeader := newHttpResponseHeader(header)

	//fmt.Println(responseHeader.description)
	conn.Write([]byte(header))

	hs.sendResponse(hs.DocRoot + requestHeader.url, *responseHeader, conn)
}

func (hs *HttpServer) sendResponse(fileName string, responseHeader HttpResponseHeader, conn net.Conn) {
	// Send headers
	// Send file if required
	// Hint - Use the bufio package to write response
	log.Println("now send data back to client...")
	file, err := os.Open(fileName)
	//log.Println()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Panicln(err)
		}
	}()

	buffer := make([]byte, 1024)
	bytesRead := 0
	for {
		size, err := file.Read(buffer)
		bytesRead += size
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		conn.Write(buffer[:size])
	}
	log.Println("total bytes sent to client: ", bytesRead)
}

func buildRespContent(version, code, description string) (content string) {
	content = ""
	content += version
	content += " "
	content += code
	content += " "
	content += description
	content += "\r\n"
	return content
}

func (hs *HttpServer) buildResponseMap(code, url string, isValid bool) (responseMap string) {
	responseMap = ""
	responseMap += "Server: "
	responseMap += " "
	responseMap += "Go-Triton-Server-1.0"
	responseMap += "\r\n"
	//if code is 200, modiefied time, type, length
	if code == "200" {
		file, err := os.Stat(hs.DocRoot + url)
		if err != nil {
			log.Println(err)
		}

		modifiedtime := file.ModTime()
		responseMap += "Last-Modified: "
		responseMap += modifiedtime.Format(time.RFC1123Z)
		responseMap += "\r\n"

		responseMap += "Content-Length:"
		log.Printf("file size is: %d\n", file.Size())
		responseMap += strconv.FormatInt(file.Size(), 10)
		log.Printf("content-length is: %s\n", strconv.FormatInt(file.Size(), 10))
		responseMap += "\r\n"

		responseMap += "Content-Type: "
		// get .xxx from url
		contentType := strings.Split(url, ".")
		responseMap += hs.MIMEMap["."+contentType[len(contentType)-1]]
		responseMap += "\r\n"
	}

	//TODO: code != 200
	responseMap += "\r\n"
	return responseMap

}
