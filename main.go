package main

import (
	"strings"
	"fmt"
	"os"
	"net"
	"reflect"
)

var MAX_REQUEST_SIZE int16 = 1024
var IP = "0.0.0.0:8000"

type Request struct {
	method string
	path string
	version string
	headers map[string]string
	body []byte
}

type Respone struct {
	status_code int
	headers map[string]string
	body []byte
}

func parse_request(request []byte) Request {
	// Handling request line
	fields := strings.Split(string(request), "\r\n") 
	method := strings.Fields(fields[0])[0]
	path := strings.Fields(fields[0])[1]
	version := strings.Fields(fields[0])[2]

	// Handling headers and body
	headers := make(map[string]string)
	var body []byte  

	for i := 1; i < len(fields); i++ {
		items := strings.SplitN(fields[i], ":", 2)
		// In case we slice white spaces
		if len(items) <= 1 { continue }
		key := items[0]
		value := items[1]

		switch key {
		case "body":
			// TO DO: create body based on content lenght
			// this is tmp
			body = []byte(value)
		default:
			headers[key] = value
		}
	}
	return Request{method: method, path: path, version: version, headers: headers, body: body}
}

func printStructFields(req Request) {
	val := reflect.ValueOf(req)
	typ := reflect.TypeOf(req)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Handle different types (e.g., []byte, map)
		switch value.Kind() {
		case reflect.Slice:
			fmt.Printf("%s: %v\n", field.Name, value.Bytes())
		case reflect.Map:
			fmt.Printf("%s:\n", field.Name)
			for _, key := range value.MapKeys() {
				fmt.Printf("  %v: %v\n", key, value.MapIndex(key))
			}
		default:
			fmt.Printf("%s: %v\n", field.Name, value)
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

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection: ", err.Error())
			os.Exit(1)
		}

		// Parsing (Check if http request is valid and parse into struct)
		// Route request (pass struct to the router and let it figure out) -> return resposne
		// sent back response

		request := make([]byte, MAX_REQUEST_SIZE)	
		conn.Read(request)
		parsed_req := parse_request(request)

		printStructFields(parsed_req)

	}
}
