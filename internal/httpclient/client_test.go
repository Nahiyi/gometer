package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	client := New(5000)

	if client.timeout != 5000*time.Millisecond {
		t.Errorf("Expected timeout 5s, got %v", client.timeout)
	}
	if client.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
}

func TestDoRequestGetSuccess(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	client := New(5000)
	result := client.DoRequest("GET", server.URL, nil, "")

	if !result.Success {
		t.Errorf("Expected success=true, got false. Error: %s", result.Error)
	}
	if result.ResponseStatus != http.StatusOK {
		t.Errorf("Expected status 200, got %d", result.ResponseStatus)
	}
	if result.ResponseTimeMs <= 0 {
		t.Errorf("ResponseTimeMs should be > 0, got %d", result.ResponseTimeMs)
	}
}

func TestDoRequestPostWithBody(t *testing.T) {
	var receivedBody string
	var receivedHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedBody = r.FormValue("data")
		receivedHeader = r.Header.Get("X-Custom-Header")
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := New(5000)
	headers := map[string]string{
		"Content-Type":    "application/x-www-form-urlencoded",
		"X-Custom-Header": "custom-value",
	}
	body := "data=test"

	result := client.DoRequest("POST", server.URL, headers, body)

	if !result.Success {
		t.Errorf("Expected success=true, got false. Error: %s", result.Error)
	}
	if receivedHeader != "custom-value" {
		t.Errorf("Expected header 'custom-value', got '%s'", receivedHeader)
	}
	if receivedBody != "test" {
		t.Errorf("Expected body 'test', got '%s'", receivedBody)
	}
}

func TestDoRequestServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	client := New(5000)
	result := client.DoRequest("GET", server.URL, nil, "")

	if result.Success {
		t.Errorf("Expected success=false for 500 error, got true")
	}
	if result.ResponseStatus != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", result.ResponseStatus)
	}
}

func TestDoRequestTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Very short timeout
	client := New(50)
	result := client.DoRequest("GET", server.URL, nil, "")

	if result.Success {
		t.Errorf("Expected success=false for timeout, got true")
	}
	if result.Error == "" {
		t.Error("Expected error message for timeout")
	}
}

func TestDoRequestInvalidURL(t *testing.T) {
	client := New(5000)
	result := client.DoRequest("GET", "http://invalid-domain-that-does-not-exist.com", nil, "")

	if result.Success {
		t.Errorf("Expected success=false for invalid URL, got true")
	}
	if result.Error == "" {
		t.Error("Expected error message for invalid URL")
	}
}

func TestDoRequestHeaderMerging(t *testing.T) {
	var contentType, auth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType = r.Header.Get("Content-Type")
		auth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(5000)
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer token123",
	}

	result := client.DoRequest("GET", server.URL, headers, "")

	if !result.Success {
		t.Errorf("Expected success=true, got false. Error: %s", result.Error)
	}
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
	if auth != "Bearer token123" {
		t.Errorf("Expected Authorization 'Bearer token123', got '%s'", auth)
	}
}
