package coreutils

import (
	"strings"
	"strconv"
	"fmt"
)
	
type Request struct {
	Method string
	Path string
	Version string
	Headers map[string]string
	Body []byte
}

type Respone struct {
	Status_code int
	Headers map[string]string
	Body []byte
}

func ParseRequest(request []byte) Request {
	raw := string(request)

	// Split into header and body
	parts := strings.SplitN(raw, "\r\n\r\n", 2)
	if len(parts) < 2 {
		return Request{}
	}
	headerPart := parts[0]
	bodyPart := parts[1]

	lines := strings.Split(headerPart, "\r\n")
	requestLine := lines[0]
	headerLines := lines[1:]

	// Parse request line
	fields := strings.Fields(requestLine)
	if len(fields) < 3 {
		return Request{}
	}
	method := fields[0]
	path := fields[1]
	version := fields[2]

	headers := make(map[string]string)
	for _, line := range headerLines {
		items := strings.SplitN(line, ":", 2)
		if len(items) != 2 {
			continue
		}
		key := strings.TrimSpace(items[0])
		value := strings.TrimSpace(items[1])
		headers[key] = value
	}

	var body []byte
	if lengthStr, ok := headers["Content-Length"]; ok {
		if length, err := strconv.Atoi(lengthStr); err == nil && length <= len(bodyPart) {
			body = []byte(bodyPart[:length])
		}
	}

	return Request{
		Method:  method,
		Path:    path,
		Version: version,
		Headers: headers,
		Body:    body,
	}
}

func FormatResponse(resp Respone) []byte{
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

	responseStr += "\r\n"

	full := append([]byte(responseStr), resp.Body...)
	return full
}
