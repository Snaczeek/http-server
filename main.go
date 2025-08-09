package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"snaczek-server/coreutils"
	"snaczek-server/router"
	"strings"
	"time"
)

var MAX_REQUEST_SIZE int16 = 1024
var IP = "0.0.0.0:8000"

func handleConnection(conn net.Conn, r *router.Router) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		request := make([]byte, MAX_REQUEST_SIZE)
		n, err := reader.Read(request)
		if err != nil {
			return // Connection close/timeout
		}

		parsedReq, err := coreutils.ParseRequest(request[:n])
		if err != nil {
			resp := coreutils.BadRequestResponse("400 Bad Request: " + err.Error())
			conn.Write(coreutils.FormatResponse(resp))
			return
		}

		fmt.Println(parsedReq.Method, parsedReq.Path)

		resp := r.Route(parsedReq)

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
