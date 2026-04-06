package reporter

import (
	"encoding/json"
	"os"
	"testing"

	"gometer/internal/collector"
)

func TestGenerateEmptyRecords(t *testing.T) {
	records := []collector.ThreadRecord{}
	report := Generate(records, 1000)

	if report.Summary.TotalThreads != 0 {
		t.Errorf("Expected 0 threads, got %d", report.Summary.TotalThreads)
	}
	if report.Summary.TotalRequests != 0 {
		t.Errorf("Expected 0 requests, got %d", report.Summary.TotalRequests)
	}
	if report.Summary.DurationMs != 1000 {
		t.Errorf("Expected duration 1000, got %d", report.Summary.DurationMs)
	}
}

func TestGenerateWithRecords(t *testing.T) {
	records := []collector.ThreadRecord{
		{
			ThreadId: 1,
			LoopResults: []collector.LoopRecord{
				{
					LoopIndex: 1,
					Requests: []collector.RequestRecord{
						{RequestIndex: 1, ResponseTimeMs: 100, Success: true},
						{RequestIndex: 2, ResponseTimeMs: 200, Success: true},
					},
				},
			},
		},
		{
			ThreadId: 2,
			LoopResults: []collector.LoopRecord{
				{
					LoopIndex: 1,
					Requests: []collector.RequestRecord{
						{RequestIndex: 1, ResponseTimeMs: 150, Success: false, Error: "server error"},
					},
				},
			},
		},
	}

	report := Generate(records, 5000)

	if report.Summary.TotalThreads != 2 {
		t.Errorf("Expected 2 threads, got %d", report.Summary.TotalThreads)
	}
	if report.Summary.TotalRequests != 3 {
		t.Errorf("Expected 3 requests, got %d", report.Summary.TotalRequests)
	}
	if report.Summary.SuccessCount != 2 {
		t.Errorf("Expected 2 successes, got %d", report.Summary.SuccessCount)
	}
	if report.Summary.FailCount != 1 {
		t.Errorf("Expected 1 failure, got %d", report.Summary.FailCount)
	}
	if report.Summary.SuccessRate != 2.0/3.0 {
		t.Errorf("Expected success rate 0.667, got %f", report.Summary.SuccessRate)
	}
	if report.Summary.MinResponseTimeMs != 100 {
		t.Errorf("Expected min 100, got %d", report.Summary.MinResponseTimeMs)
	}
	if report.Summary.MaxResponseTimeMs != 200 {
		t.Errorf("Expected max 200, got %d", report.Summary.MaxResponseTimeMs)
	}
	if report.Summary.AvgResponseTimeMs != 150.0 {
		t.Errorf("Expected avg 150, got %f", report.Summary.AvgResponseTimeMs)
	}
}

func TestPercentile(t *testing.T) {
	// 100 elements: 1, 2, 3, ..., 100
	sorted := make([]int64, 100)
	for i := 0; i < 100; i++ {
		sorted[i] = int64(i + 1)
	}

	if p50 := percentile(sorted, 50); p50 != 50 {
		t.Errorf("Expected P50=50, got %d", p50)
	}
	if p90 := percentile(sorted, 90); p90 != 90 {
		t.Errorf("Expected P90=90, got %d", p90)
	}
	if p99 := percentile(sorted, 99); p99 != 99 {
		t.Errorf("Expected P99=99, got %d", p99)
	}
}

func TestPercentileEmpty(t *testing.T) {
	if percentile([]int64{}, 50) != 0 {
		t.Error("Expected 0 for empty slice")
	}
}

func TestWriteToStdout(t *testing.T) {
	report := Report{
		Summary: Summary{TotalRequests: 10},
	}

	// This should just print to stdout without error
	err := Write(report, "")
	if err != nil {
		t.Errorf("Write to stdout failed: %v", err)
	}
}

func TestWriteToFile(t *testing.T) {
	report := Report{
		Summary: Summary{TotalRequests: 10},
	}

	tmpfile, err := os.CreateTemp("", "report-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	err = Write(report, tmpfile.Name())
	if err != nil {
		t.Errorf("Write to file failed: %v", err)
	}

	// Verify file content
	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Errorf("Read back file failed: %v", err)
	}

	var loaded Report
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Errorf("Unmarshal failed: %v", err)
	}
	if loaded.Summary.TotalRequests != 10 {
		t.Errorf("Expected 10 requests in file, got %d", loaded.Summary.TotalRequests)
	}
}
