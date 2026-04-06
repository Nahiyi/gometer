package collector

import (
	"sync"
)

// RequestRecord represents a single request result
type RequestRecord struct {
	RequestIndex    int               `json:"requestIndex"`
	URL             string            `json:"url"`
	Method          string            `json:"method"`
	RequestHeaders  map[string]string `json:"requestHeaders"`
	RequestBody     string            `json:"requestBody"`
	ResponseStatus  int               `json:"responseStatus"`
	ResponseTimeMs  int64             `json:"responseTimeMs"`
	ResponseHeaders map[string]string `json:"responseHeaders"`
	Success         bool              `json:"success"`
	Error           string            `json:"error"`
}

// LoopRecord represents results of one loop iteration
type LoopRecord struct {
	LoopIndex int             `json:"loopIndex"`
	Requests  []RequestRecord `json:"requests"`
}

// ThreadRecord represents all results for one thread
type ThreadRecord struct {
	ThreadId    int           `json:"threadId"`
	LoopResults []LoopRecord  `json:"loopResults"`
}

// Collector collects results from all threads in a thread-safe manner
type Collector struct {
	mu           sync.Mutex
	threadRecords []ThreadRecord
}

// New creates a new Collector
func New() *Collector {
	return &Collector{
		threadRecords: make([]ThreadRecord, 0),
	}
}

// NewThreadRecord initializes a new thread record and returns its index
func (c *Collector) NewThreadRecord(threadId int) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.threadRecords = append(c.threadRecords, ThreadRecord{
		ThreadId:    threadId,
		LoopResults: make([]LoopRecord, 0),
	})
	return len(c.threadRecords) - 1
}

// AddLoopResult adds a loop result to a thread
func (c *Collector) AddLoopResult(threadIndex, loopIndex int, requests []RequestRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	loopRecord := LoopRecord{
		LoopIndex: loopIndex,
		Requests:  requests,
	}
	c.threadRecords[threadIndex].LoopResults = append(c.threadRecords[threadIndex].LoopResults, loopRecord)
}

// GetAllRecords returns all collected thread records (copy)
func (c *Collector) GetAllRecords() []ThreadRecord {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := make([]ThreadRecord, len(c.threadRecords))
	copy(result, c.threadRecords)
	return result
}
