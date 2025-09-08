# Snaczek HTTP Server

A lightweight HTTP/1.1 server written in **Go** from scratch.  
Ive built this project to understand more about networking, HTTP protocol internals, and low-level server implementation. 

---

## Features

- **HTTP/1.1 compliant basics**
  - Request parsing (method, path, version, headers, body)
  - Response formatting with proper headers
  - Persistent connections (`keep-alive` / `close`)
  - Host header validation
- **Routing system**
  - Register routes with custom handlers (Similar to django views)
  - Support for multiple HTTP methods (GET, POST, DELETE, etc.)
- **Request body parsing**
  - Content-Length support
  - JSON body parsing 
- **Response helpers**
  - Error response templates (400, 404, etc.)
  - Gzip compression support (`Accept-Encoding: gzip`)
- **Concurrency**
  - Handles multiple clients concurrently with goroutines
  - Per-connection 5s timeout

---

## How to Run

### 1. Dependencies
- Go version **1.24.5** (or higher)

### 2. Run the server
```bash
go run .
```

By default the server listens on 0.0.0.0:8000.
You can change the address/port by passing it as an argument:
```bash 
go run . 127.0.0.1:9000
```

---

## How to Create Routes & Handlers

Handlers are just Go functions that:

- Receive a parsed request (coreutils.Request)
- Return a response (coreutils.Response)

### 1. Define a handler
```go
package handlers

import "snaczek-server/coreutils"

func HelloHandler(req coreutils.Request) coreutils.Response {
	return coreutils.Response{
		Status_code: 200,
		Headers: map[string]string{"Content-Type": "text/plain"},
		Body: []byte("Hello from GET /hello\n"),
	}
}
```

### 2. Register the route

Inside your `routes.go`, register the handler with the router:
```go
func RegisterAllRoutes(r *router.Router) {
	r.RegisterRoute("GET", "/hello", handlers.HelloHandler)
}
```
Now when you run the server, and visit:
```bash
curl http://0.0.0.0:8000/hello
```
You will get: 
```bash
Hello from GET /hello
```

---

## Request & Response Structs

### coreutils.Request

```go
type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    []byte
}

```

### coreutils.Response
```go
type Response struct {
	Status_code int
	Headers     map[string]string
	Body        []byte
}

```
