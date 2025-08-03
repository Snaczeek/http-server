package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"snaczek-server/coreutils"
	"snaczek-server/router"
	"strings"
)

var MAX_REQUEST_SIZE int16 = 1024
var IP = "0.0.0.0:8000"

func parse_request(request []byte) coreutils.Request {
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
	return coreutils.Request{Method: method, Path: path, Version: version, Headers: headers, Body: body}
}

func format_response(resp coreutils.Respone) []byte{
	statusText := map[int]string{
		200: "OK",
		201: "Created",
		400: "Bad Request",
		404: "Not Found",
		405: "Method Not Allowed",
		500: "Internal Server Error",
	}

	statusMsg, exist := statusText[resp.Status_code]
	if !exist {
		statusMsg = "Unknow"
	}

	responseStr := fmt.Sprintf("HTTP/1.1 %d %s\r\n", resp.Status_code, statusMsg)

	if resp.Headers == nil {
		resp.Headers = make(map[string]string)
	}
	resp.Headers["Content-Length"] = fmt.Sprintf("%d", len(resp.Body))

	for key, value := range resp.Headers {
		responseStr += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	full := append([]byte(responseStr), resp.Body...)
	return full
}


func printStructFields(req coreutils.Request) {
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

	r := router.NewRouter()
	r.RegisterRoute("GET", "/hello", func(req coreutils.Request) coreutils.Respone {
		return coreutils.Respone{
			Status_code: 200,
			Headers: map[string]string{"Content-Type": "text/plain"},
			Body: []byte("Hello from GET /hello"),
		}
	})

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


		response := r.Route(parsed_req)
		raw := format_response(response)
		conn.Write(raw)
	}
}
