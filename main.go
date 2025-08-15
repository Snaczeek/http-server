package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"snaczek-server/coreutils"
	"snaczek-server/router"
	"strconv"
	"strings"
	"time"
	"compress/gzip"
)

var IP = "0.0.0.0:8000"
var COMPRESSION_THRESHOLD = 500 // Default minimum 500 bytes  

func handleConnection(conn net.Conn, r *router.Router) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		request, err := readHTTPRequest(reader)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded){
				fmt.Println("Connection timeout")
			} else if err == io.EOF {
				fmt.Println("Client closed connection")
			} else {
				fmt.Print("Read error:", err)
			}
			return
		}

		parsedReq, err := coreutils.ParseRequest(request)
		if err != nil {
			resp := coreutils.BadRequestResponse("400 Bad Request: " + err.Error())
			conn.Write(coreutils.FormatResponse(resp))
			return
		}

		fmt.Println(parsedReq.Method, parsedReq.Path)

		resp := r.Route(parsedReq)
		acceptEnc := parsedReq.Headers["Accept-Encoding"]

		// Compression handling
		if strings.Contains(acceptEnc, "gzip") && len(resp.Body) >= COMPRESSION_THRESHOLD{
			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			gz.Write(resp.Body)
			gz.Close()

			resp.Body = buf.Bytes()
			resp.Headers["Content-Encoding"] = "gzip"
		}

		// Persistent connection handling
		connHeader := strings.ToLower(parsedReq.Headers["Connection"])
		if connHeader == "close" {
			resp.Headers["Connection"] = "close"
			return 
		} else {
			resp.Headers["Connection"] = "keep-alive"
			conn.Write(coreutils.FormatResponse(resp))
		}
	}
}

func readHTTPRequest(reader *bufio.Reader) ([]byte, error) {
	var buffer bytes.Buffer

	// Read headers
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		buffer.WriteString(line)
		if line == "\r\n" {
			break
		}
	}

	// Handle body
	headers := parseTmpHeaders(buffer.Bytes())
	if clStr, ok := headers["Content-Length"]; ok {
		cl, err := strconv.Atoi(clStr)
		if err != nil {
			return nil, fmt.Errorf("invalid Content-Length")
		}
		body := make([]byte, cl)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			return nil, err
		}
		buffer.Write(body)
	}

	return buffer.Bytes(), nil
}

func parseTmpHeaders(raw []byte) map[string]string {
	headers := make(map[string]string)
	lines := strings.Split(string(raw), "\r\n")
	for _, line := range lines[1:] { // Skip request line
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return headers
}

func main () {
	main_args := os.Args
	if len(main_args) > 1 { IP = main_args[1] }

	fmt.Println("Starting server at", IP)
	ln, err := net.Listen("tcp", IP)
	if err != nil {
		fmt.Println("Failed to bind", IP)
		os.Exit(1)
	}

	r := router.NewRouter()
	router.RegisterAllRoutes(r)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection: ", err.Error())
			continue
		}

		go handleConnection(conn, r)
	}
}
