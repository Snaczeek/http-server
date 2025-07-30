package main

import (
	"fmt"
	"os"
	"net"
	"strconv"
)

var MAX_REQUEST_SIZE int64 = 1024
var IP = "0.0.0.0:8000"

type Request struct {
	method string
	path string
	headers map[string]string
	body []byte
}

type Respone struct {
	status_code int
	headers map[string]string
	body []byte
}

func parse_request(request []byte) Request {
	return Request{}
}

func main () {
	main_args := os.Args
	if len(main_args) > 1 { IP = main_args[1] }
	if len(main_args) > 2 { 
		tmp, err := strconv.ParseInt(main_args[2], 10, 16)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		MAX_REQUEST_SIZE = tmp
	}

	fmt.Println("Starting server at", IP)
	ln, err := net.Listen("tcp", IP)
	if err != nil {
		fmt.Println("Failed to bind", IP)
		os.Exit(1)
	}

	for {
		_, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection: ", err.Error())
			os.Exit(1)
		}

		// Parsing (Check if http request is valid and parse into struct)
		// Route request (pass struct to the router and let it figure out) -> return resposne
		// sent back response
	}
}
