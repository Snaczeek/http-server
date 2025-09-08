package handlers

import (
	"snaczek-server/coreutils"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func CreateUserHandler(req coreutils.Request) coreutils.Response {
	user, err := coreutils.ParseJSONBody[User](req)
	if err != nil {
		return coreutils.Response{
			Status_code: 400,
			Headers: map[string]string{"Content-Type": "application/json"},
			Body: []byte(`{"error":"invalid JSON"}`),
		}
	}

	responseBody := []byte(`{"message":"User created: ` + user.Name + `"}`)

	return coreutils.Response{
		Status_code: 201,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: responseBody,
	}
}
