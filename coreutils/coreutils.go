package coreutils

import (
	"strings"
	"strconv"
	"fmt"
	"net"
	"time"
	"encoding/json"
	"errors"
)
	
type Request struct {
	Method string
	Path string
	Version string
	Headers map[string]string
	Body []byte
}

type Response struct {
	Status_code int
	Headers map[string]string
	Body []byte
}

func isValidHost(host string) bool {
	// Split into host and optional port
	h, p, err := net.SplitHostPort(host)
	if err != nil {
		// No port case
		h = host
		p = ""
	}

	// Validate hostname or IP
	if net.ParseIP(h) == nil {
		if !isValidHostname(h) {
			return false
		}
	}

	// Validate port if present
	if p != "" {
		if _, err := strconv.Atoi(p); err != nil {
			return false
		}
	}

	return true
}

func isValidHostname(h string) bool {
	// Hostname must be <= 253 chars and split into labels <= 63 chars
	if len(h) == 0 || len(h) > 253 {
		return false
	}
	labels := strings.Split(h, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return false
		}
		// Only letters, digits, hyphen, and can't start/end with hyphen
		for i, c := range label {
			if !(c == '-' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9') {
				return false
			}
			if (i == 0 || i == len(label)-1) && c == '-' {
				return false
			}
		}
	}
	return true
}

func ParseRequest(request []byte) (Request, error) {
	raw := string(request)

	// Split into header and body
	parts := strings.SplitN(raw, "\r\n\r\n", 2)
	if len(parts) < 2 {
		return Request{}, fmt.Errorf("Invalid request: no headers/body seperation")
	}
	headerPart := parts[0]
	bodyPart := parts[1]

	lines := strings.Split(headerPart, "\r\n")
	requestLine := lines[0]
	headerLines := lines[1:]

	// Parse request line
	fields := strings.Fields(requestLine)
	if len(fields) < 3 {
		return Request{}, fmt.Errorf("Invalid reqeust line")
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

	// HTTP/1.1 Host header check
	if version == "HTTP/1.1" {
		host, ok := headers["Host"]
		if !ok || !isValidHost(host) {
			return Request{}, fmt.Errorf("invalid or missing host name")
		}
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
	}, nil
}

func FormatResponse(resp Response) []byte{
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
	curr_date := time.Now().UTC().Format(time.RFC1123)
	curr_date = fmt.Sprintf("%sGMT", curr_date[:len(curr_date)-3])
	resp.Headers["Date"] = curr_date

	for key, value := range resp.Headers {
		responseStr += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	responseStr += "\r\n"

	full := append([]byte(responseStr), resp.Body...)
	return full
}

func BadRequestResponse(message string) Response {
	body := []byte(message + "\n")
	return Response{
		Status_code: 400,
		Headers: map[string]string{
			"Content-Type":   "text/plain; charset=utf-8",
			"Content-Length": fmt.Sprintf("%d", len(body)),
			"Connection":     "close", 
		},
		Body: body,
	}
}

func ParseJSONBody[T any](req Request) (T, error) {
	var data T
	contentType := req.Headers["Content-Type"]

	if contentType != "application/json" {
		return data, errors.New("unsupported media type")
	}

	err := json.Unmarshal(req.Body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
