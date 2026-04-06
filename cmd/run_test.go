package cmd

import (
	"os"
	"testing"
)

func TestDryRunWithValidConfig(t *testing.T) {
	// Create a temp config file
	content := `{
		"request": {
			"url": "http://example.com/api",
			"method": "POST",
			"headers": {"Content-Type": "application/json"},
			"body": "{}"
		},
		"users": [
			{"headers": {"token": "user-token-1"}}
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

	// Set up flags for dry-run
	configFile = tmpfile.Name()
	dryRun = true

	// Run dry-run
	err = runPressureTest(nil, nil)
	if err != nil {
		t.Errorf("dry-run should not fail: %v", err)
	}
}

func TestDryRunWithInvalidJSON(t *testing.T) {
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

	configFile = tmpfile.Name()
	dryRun = true

	err = runPressureTest(nil, nil)
	if err == nil {
		t.Error("dry-run with invalid JSON should fail")
	}
}

func TestValidateThreadCount(t *testing.T) {
	content := `{
		"request": {
			"url": "http://example.com/api",
			"method": "GET"
		},
		"users": [{"headers": {"token": "user1"}}]
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

	configFile = tmpfile.Name()
	threads = -1 // Invalid
	dryRun = false

	err = runPressureTest(nil, nil)
	if err == nil {
		t.Error("negative threads should fail")
	}
	if err != nil && err.Error() != "threads must be > 0" {
		t.Errorf("wrong error message: %v", err)
	}
}

func TestValidateUserCount(t *testing.T) {
	content := `{
		"request": {
			"url": "http://example.com/api",
			"method": "GET"
		},
		"users": [{"headers": {"token": "user1"}}]
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

	configFile = tmpfile.Name()
	threads = 10 // More than users
	dryRun = false

	err = runPressureTest(nil, nil)
	if err == nil {
		t.Error("insufficient users should fail")
	}
}
