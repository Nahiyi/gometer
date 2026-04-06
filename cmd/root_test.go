package cmd

import (
	"testing"
)

func TestFlagDefinitions(t *testing.T) {
	// Test that all flags are properly defined
	flags := runCmd.Flags()

	// threads flag
	if flags.Lookup("threads") == nil {
		t.Error("threads flag should exist")
	}

	// ramp-up flag
	rampUpFlag := flags.Lookup("ramp-up")
	if rampUpFlag == nil {
		t.Error("ramp-up flag should exist")
	}
	if rampUpFlag.DefValue != "0" {
		t.Errorf("ramp-up default should be 0, got %s", rampUpFlag.DefValue)
	}

	// loop flag
	loopFlag := flags.Lookup("loop")
	if loopFlag == nil {
		t.Error("loop flag should exist")
	}
	if loopFlag.DefValue != "1" {
		t.Errorf("loop default should be 1, got %s", loopFlag.DefValue)
	}

	// config flag
	configFlag := flags.Lookup("config")
	if configFlag == nil {
		t.Error("config flag should exist")
	}
	if configFlag.DefValue != "./req.json" {
		t.Errorf("config default should be ./req.json, got %s", configFlag.DefValue)
	}

	// request-timeout flag
	timeoutFlag := flags.Lookup("request-timeout")
	if timeoutFlag == nil {
		t.Error("request-timeout flag should exist")
	}
	if timeoutFlag.DefValue != "5000" {
		t.Errorf("request-timeout default should be 5000, got %s", timeoutFlag.DefValue)
	}

	// max-duration flag
	maxDurFlag := flags.Lookup("max-duration")
	if maxDurFlag == nil {
		t.Error("max-duration flag should exist")
	}
	if maxDurFlag.DefValue != "0" {
		t.Errorf("max-duration default should be 0, got %s", maxDurFlag.DefValue)
	}

	// dry-run flag
	dryRunFlag := flags.Lookup("dry-run")
	if dryRunFlag == nil {
		t.Error("dry-run flag should exist")
	}
}
