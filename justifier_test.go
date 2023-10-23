package headersvalidator_test

import (
	"context"
	"github.com/Farrukhraz/headersvalidator"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithHeaders(t *testing.T) {
	key := "X-AUTH-TOKEN"
	value := "12345"

	cfg := headersvalidator.CreateConfig()
	cfg.Headers["key"] = key
	cfg.Headers["value"] = value

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := headersvalidator.New(ctx, next, cfg, "demo-plugin")
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://google.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(key, value)

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, key, value)
}

func TestWithoutHeaders(t *testing.T) {
	key := "X-AUTH-TOKEN"
	value := "12345"
	cfg := headersvalidator.CreateConfig()
	cfg.Headers["key"] = key
	cfg.Headers["value"] = value

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := headersvalidator.New(ctx, next, cfg, "demo-plugin")
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://google.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(recorder, req)

	if recorder.Code != 401 {
		t.Errorf("401 response was expected, but doesn't received")
	}
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", req.Header.Get(key))
	}
}
