package main

import (
	"encoding/json"
	"log"
	"net"
	"strings"
)

/*
ALGO:
- Listen for TCP connections
- Accept connection
- Process it
*/

type Request struct {
	RequestLine string
	Headers     map[string]string
	// Body        map[string]string
}

func logStruct(v any) {
	json, _ := json.MarshalIndent(v, "", "  ")
	log.Printf("%+v\n", string(json))
}

func getHeaders(requestRaw string) (headers map[string]string) {
	requestLines := strings.Split(requestRaw, "\n")
	if len(requestLines) <= 1 {
		return map[string]string{}
	}

	i := 1
	headers = map[string]string{}
	for {
		line := requestLines[i]
		if len(strings.TrimSpace(line)) == 0 || i == len(requestLines)-1 {
			break
		}
		header := strings.Split(line, ": ")
		headerName := header[0]
		headerValue := header[1]
		headers[headerName] = headerValue
		i++
	}

	return headers
}

func main() {
	server, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		log.Fatalln(err)
	}

	defer server.Close()

	buffer := make([]byte, 4000)
	for {
		connection, err := server.Accept()
		log.Println("New connection established")
		if err != nil {
			log.Fatalln(err)
		}

		byteCount, _ := connection.Read(buffer)
		requestRaw := strings.ReplaceAll(string(buffer), "\r", "")
		requestLines := strings.Split(requestRaw, "\n")
		if len(requestLines) == 0 {
			log.Println("Invalid HTTP request")
		}
		log.Println(byteCount)

		// logStruct(requestLines)
		// requestLine := requestLines[0]

		headers := getHeaders(requestRaw)

		requestInfo := Request{
			RequestLine: requestLines[0],
			Headers:     headers,
		}

		logStruct(requestInfo)
		responseBytes, _ := json.Marshal(requestInfo)
		responseLine := "HTTP/1.1 200 OK\r\n\r\n"
		connection.Write([]byte(responseLine))
		connection.Write(responseBytes)

		connection.Close() // Try to not have this
	}
}
