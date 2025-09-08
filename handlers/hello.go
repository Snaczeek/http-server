package handlers

import "snaczek-server/coreutils"

func HelloHandler(req coreutils.Request) coreutils.Response {
	return coreutils.Response{
		Status_code: 200,
		Headers: map[string]string{"Content-Type": "text/plain"},
		Body: []byte("Hello from GET /hello\n"),
	}
}
