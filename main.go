package main

import (
	"fmt"
	"net"
	"os"
	"snaczek-server/coreutils"
	"snaczek-server/router"
)

var MAX_REQUEST_SIZE int16 = 1024
var IP = "0.0.0.0:8000"

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
			os.Exit(1)
		}

		request := make([]byte, MAX_REQUEST_SIZE)	
		conn.Read(request)
		parsed_req, err := coreutils.ParseRequest(request)
		if err != nil {
			response := coreutils.BadRequestResponse("400 Bad Request: " + err.Error())
			raw := coreutils.FormatResponse(response)
			conn.Write(raw)
		}
		fmt.Println(parsed_req.Path + " " + parsed_req.Method)

		response := r.Route(parsed_req)
		raw := coreutils.FormatResponse(response)
		conn.Write(raw)
	}
}
