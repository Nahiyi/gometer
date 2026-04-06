package config

import (
	"os"
	"testing"
)

func TestLoadValidConfig(t *testing.T) {
	// Create a temp config file
	content := `{
		"request": {
			"url": "http://example.com/api",
			"method": "POST",
			"headers": {
				"Content-Type": "application/json"
			},
			"body": "{\"key\": \"value\"}"
		},
		"users": [
			{
				"headers": {
					"token": "user-token-1",
					"x-user-id": "1001"
				}
			},
			{
				"headers": {
					"token": "user-token-2",
					"x-user-id": "1002"
				}
			}
		]
	}`

	tmpfile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify request
	if cfg.Request.URL != "http://example.com/api" {
		t.Errorf("URL expected 'http://example.com/api', got '%s'", cfg.Request.URL)
	}
	if cfg.Request.Method != "POST" {
		t.Errorf("Method expected 'POST', got '%s'", cfg.Request.Method)
	}
	if cfg.Request.Headers["Content-Type"] != "application/json" {
		t.Errorf("Content-Type header incorrect")
	}
	if cfg.Request.Body != "{\"key\": \"value\"}" {
		t.Errorf("Body incorrect")
	}

	// Verify users
	if len(cfg.Users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(cfg.Users))
	}
	if cfg.Users[0].Headers["token"] != "user-token-1" {
		t.Errorf("First user token incorrect")
	}
	if cfg.Users[1].Headers["x-user-id"] != "1002" {
		t.Errorf("Second user x-user-id incorrect")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte("invalid json")); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = Load(tmpfile.Name())
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.json")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
