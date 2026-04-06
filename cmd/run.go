package cmd

import (
	"fmt"
	"sync"
	"time"

	"gometer/internal/collector"
	"gometer/internal/config"
	"gometer/internal/httpclient"
	"gometer/internal/loader"
	"gometer/internal/reporter"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run pressure test",
	RunE:  runPressureTest,
}

func runPressureTest(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Dry run mode - just validate config
	if dryRun {
		fmt.Println("Config is valid:")
		fmt.Printf("  URL: %s\n", cfg.Request.URL)
		fmt.Printf("  Method: %s\n", cfg.Request.Method)
		fmt.Printf("  Users: %d\n", len(cfg.Users))
		return nil
	}

	// Validate threads
	if threads <= 0 {
		return fmt.Errorf("threads must be > 0")
	}

	// Validate user count
	if len(cfg.Users) < threads {
		return fmt.Errorf("insufficient user configs: %d users for %d threads", len(cfg.Users), threads)
	}

	// Create HTTP client
	client := httpclient.New(requestTimeout)

	// Create loader
	userLoader := loader.New(cfg.Users)

	// Create collector
	coll := collector.New()

	// Calculate delay between thread starts
	var startDelay time.Duration
	if rampUp > 0 && threads > 0 {
		startDelay = time.Duration(rampUp) * time.Second / time.Duration(threads)
	}

	startTime := time.Now()
	stopCh := make(chan struct{})

	// Start max duration watcher if set
	if maxDuration > 0 {
		go func() {
			time.AfterFunc(time.Duration(maxDuration)*time.Second, func() {
				close(stopCh)
			})
		}()
	}

	// WaitGroup for thread coordination
	var wg sync.WaitGroup

	// Launch threads
	for i := 0; i < threads; i++ {
		select {
		case <-stopCh:
			// Max duration reached, stop launching new threads
			goto afterLoop
		default:
		}

		threadIndex := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			runThread(threadIndex, cfg, client, userLoader, coll, loop)
		}()

		if startDelay > 0 {
			time.Sleep(startDelay)
		}
	}
afterLoop:

	wg.Wait()
	duration := time.Since(startTime)

	// Generate report
	records := coll.GetAllRecords()
	report := reporter.Generate(records, duration.Milliseconds())

	return reporter.Write(report, outputFile)
}

func runThread(threadIndex int, cfg *config.Config, client *httpclient.Client, userLoader *loader.Loader, coll *collector.Collector, loopCount int) {
	idx := coll.NewThreadRecord(threadIndex)
	userConfig := userLoader.GetUserConfig(threadIndex)

	for loopIndex := 1; loopIndex <= loopCount; loopIndex++ {
		// Build merged headers (shared + user-specific)
		mergedHeaders := make(map[string]string)
		for k, v := range cfg.Request.Headers {
			mergedHeaders[k] = v
		}
		for k, v := range userConfig.Headers {
			mergedHeaders[k] = v
		}

		// Execute request
		result := client.DoRequest(
			cfg.Request.Method,
			cfg.Request.URL,
			mergedHeaders,
			cfg.Request.Body,
		)

		reqRecord := collector.RequestRecord{
			RequestIndex:    1,
			URL:             cfg.Request.URL,
			Method:          cfg.Request.Method,
			RequestHeaders:  mergedHeaders,
			RequestBody:     cfg.Request.Body,
			ResponseStatus:  result.ResponseStatus,
			ResponseTimeMs:  result.ResponseTimeMs,
			ResponseHeaders: result.ResponseHeaders,
			Success:         result.Success,
			Error:           result.Error,
		}

		coll.AddLoopResult(idx, loopIndex, []collector.RequestRecord{reqRecord})
	}
}
