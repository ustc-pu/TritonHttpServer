package tritonhttp

import (
	"errors"
	"log"
	"net"
	"strings"
)

/** 
	Initialize the tritonhttp server by populating HttpServer structure
**/
func NewHttpdServer(port, docRoot, mimePath string) (*HttpServer, error) {
	//panic("todo - NewHttpdServer")
	log.Println("building a new http server")
	// Initialize mimeMap for server to refer
	var server HttpServer
	server.ServerPort = port
	server.DocRoot = docRoot
	server.MIMEPath = mimePath
	var err error
	server.MIMEMap, err =  ParseMIME(mimePath)

	// Return pointer to HttpServer
	return &server, err
}

/** 
	Start the tritonhttp server
**/
func (hs *HttpServer) Start() (err error) {
	//panic("todo - StartServer")

	// Start listening to the server port
	l, err := net.Listen("tcp", "localhost"+ hs.ServerPort)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Listening to connections at localhost on port", hs.ServerPort)
	defer l.Close()
	// Accept connection from client

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}
		go hs.handleConnection(conn)
	}

	// Spawn a go routine to handle request
	return err
}

/**
	Initialize the tritonhttp request header by populating HttpRequestHeader structure
**/
func newHttpRequestHeader(header string) (*HttpRequestHeader, error) {
	//panic("todo - NewHttpdServer")
	log.Println("building a new http request header")
	var err error
	var requestHeader HttpRequestHeader
	requestHeader.requestHeaderMap = make(map[string]string)
	i := 0
	for strings.Contains(header, "\r\n") {
		idx := strings.Index(header, "\r\n")
		//log.Printf("index is %d\n", idx)
		//log.Println(header)
		content := []string{}
		if len(header[:idx]) == 0 {
			//log.Println("no remaining")
			break
		}
		if i == 0 {
			//add three fields from first line
			content = strings.Split(header[:idx], " ")
			requestHeader.method = content[0]
			//check default dir and default subdir
			length := len(content[1])
			if content[1][length-1] == '/' {
				requestHeader.url = content[1] + "index.html"
			} else if content[1][length-7:] == "subdir1" || content[1][length-8:] == "subdir11" {
				requestHeader.url = content[1] + "/index.html"
			} else {
				requestHeader.url = content[1]
			}
			if len(content) < 3 {
				requestHeader.version = " "
			} else {
				requestHeader.version = content[2]
			}
			i++
		} else {
			//add kv pair to map
			content = strings.Split(header[:idx], ":")
			if !strings.Contains(header, ":") {
				log.Println("malformed request header, no : in map")
				err = errors.New("malformed request header")
				return &requestHeader, err
			} else {
				requestHeader.requestHeaderMap[content[0]] = content[1]
			}
		}
		header = header[idx+2:]

		if err != nil {
			log.Println(err)
		}
	}
	return &requestHeader, err
}

func newHttpResponseHeader(content string) (*HttpResponseHeader) {
	var responseHeader HttpResponseHeader
	//TODO: initialize data in
	return &responseHeader
}

