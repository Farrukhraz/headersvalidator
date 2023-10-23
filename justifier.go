package headersvalidator

import (
	"context"
	"fmt"
	"net/http"
)

type Config struct {
	Headers map[string]string `json:"headers,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		Headers: make(map[string]string),
	}
}

type Demo struct {
	next    http.Handler
	headers map[string]string
	name    string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	key := config.Headers["key"]
	value := config.Headers["value"]
	if key == "" || value == "" {
		return nil, fmt.Errorf("'key:value' pair cannot be empty")
	}

	return &Demo{
		headers: config.Headers,
		next:    next,
		name:    name,
	}, nil
}

func (a *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// TODO задавать x-auth-token можно из docker-compose:labels.
	//  в таком случае каждый сервис сам сможет задавать для себя нужные токены
	key := a.headers["key"]
	value := a.headers["value"]
	requestHeader := req.Header.Get(key)
	if requestHeader != value {
		errorMessage := "Required authentication header is not found or invalid token is given"
		fmt.Println(errorMessage, req)
		http.Error(rw, errorMessage, http.StatusUnauthorized)
		return
	}

	a.next.ServeHTTP(rw, req)
}
