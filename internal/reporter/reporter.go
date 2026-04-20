package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"gmeter/internal/collector"
)

// Summary represents aggregated statistics
type Summary struct {
	TotalThreads      int     `json:"totalThreads"`
	TotalLoops        int     `json:"totalLoops"`
	TotalRequests     int     `json:"totalRequests"`
	SuccessCount      int     `json:"successCount"`
	FailCount         int     `json:"failCount"`
	SuccessRate       float64 `json:"successRate"`
	DurationMs        int64   `json:"durationMs"`
	AvgResponseTimeMs float64 `json:"avgResponseTimeMs"`
	MinResponseTimeMs int64   `json:"minResponseTimeMs"`
	MaxResponseTimeMs int64   `json:"maxResponseTimeMs"`
	P50ResponseTimeMs int64   `json:"p50ResponseTimeMs"`
	P90ResponseTimeMs int64   `json:"p90ResponseTimeMs"`
	P99ResponseTimeMs int64   `json:"p99ResponseTimeMs"`
}

// Report represents the full JSON report
type Report struct {
	Summary Summary                  `json:"summary"`
	Threads []collector.ThreadRecord `json:"threads"`
}

// Generate creates a report from collected results
func Generate(records []collector.ThreadRecord, durationMs int64) Report {
	var totalRequests, successCount, failCount, totalLoops int
	var totalResponseTime int64
	var minTime, maxTime int64 = -1, 0
	var allResponseTimes []int64

	for _, thread := range records {
		totalLoops += len(thread.LoopResults)
		for _, loop := range thread.LoopResults {
			for _, req := range loop.Requests {
				totalRequests++
				if req.Success {
					successCount++
				} else {
					failCount++
				}
				totalResponseTime += req.ResponseTimeMs
				allResponseTimes = append(allResponseTimes, req.ResponseTimeMs)
				if minTime == -1 || req.ResponseTimeMs < minTime {
					minTime = req.ResponseTimeMs
				}
				if req.ResponseTimeMs > maxTime {
					maxTime = req.ResponseTimeMs
				}
			}
		}
	}

	sort.Slice(allResponseTimes, func(i, j int) bool { return allResponseTimes[i] < allResponseTimes[j] })

	var avgTime float64
	if totalRequests > 0 {
		avgTime = float64(totalResponseTime) / float64(totalRequests)
	}

	summary := Summary{
		TotalThreads:      len(records),
		TotalLoops:        totalLoops,
		TotalRequests:     totalRequests,
		SuccessCount:      successCount,
		FailCount:         failCount,
		SuccessRate:       0,
		DurationMs:        durationMs,
		AvgResponseTimeMs: avgTime,
		MinResponseTimeMs: minTime,
		MaxResponseTimeMs: maxTime,
		P50ResponseTimeMs: percentile(allResponseTimes, 50),
		P90ResponseTimeMs: percentile(allResponseTimes, 90),
		P99ResponseTimeMs: percentile(allResponseTimes, 99),
	}

	if totalRequests > 0 {
		summary.SuccessRate = float64(successCount) / float64(totalRequests)
	}

	return Report{
		Summary: summary,
		Threads: records,
	}
}

func percentile(sorted []int64, p int) int64 {
	if len(sorted) == 0 {
		return 0
	}
	index := (len(sorted)-1)*p/100 + 1
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	return sorted[index-1]
}

// Write writes the report to a file or stdout
func Write(report Report, outputPath string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %v", err)
	}

	if outputPath == "" {
		fmt.Println(string(data))
		return nil
	}

	return os.WriteFile(outputPath, data, 0644)
}
