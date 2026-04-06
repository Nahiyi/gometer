package collector

import (
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	c := New()

	if c == nil {
		t.Fatal("New() returned nil")
	}
	if c.threadRecords == nil {
		t.Error("threadRecords should be initialized")
	}
}

func TestNewThreadRecord(t *testing.T) {
	c := New()

	idx1 := c.NewThreadRecord(1)
	idx2 := c.NewThreadRecord(2)

	if idx1 != 0 {
		t.Errorf("Expected first index 0, got %d", idx1)
	}
	if idx2 != 1 {
		t.Errorf("Expected second index 1, got %d", idx2)
	}

	records := c.GetAllRecords()
	if len(records) != 2 {
		t.Errorf("Expected 2 thread records, got %d", len(records))
	}
	if records[0].ThreadId != 1 {
		t.Errorf("Expected threadId 1, got %d", records[0].ThreadId)
	}
	if records[1].ThreadId != 2 {
		t.Errorf("Expected threadId 2, got %d", records[1].ThreadId)
	}
}

func TestAddLoopResult(t *testing.T) {
	c := New()

	idx := c.NewThreadRecord(1)

	requests := []RequestRecord{
		{RequestIndex: 1, URL: "http://example.com", Success: true},
	}
	c.AddLoopResult(idx, 1, requests)

	records := c.GetAllRecords()
	if len(records) != 1 {
		t.Fatalf("Expected 1 thread record, got %d", len(records))
	}
	if len(records[0].LoopResults) != 1 {
		t.Fatalf("Expected 1 loop result, got %d", len(records[0].LoopResults))
	}
	if records[0].LoopResults[0].LoopIndex != 1 {
		t.Errorf("Expected loopIndex 1, got %d", records[0].LoopResults[0].LoopIndex)
	}
	if len(records[0].LoopResults[0].Requests) != 1 {
		t.Errorf("Expected 1 request, got %d", len(records[0].LoopResults[0].Requests))
	}
}

func TestConcurrentAccess(t *testing.T) {
	c := New()

	var wg sync.WaitGroup
	numThreads := 100

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(threadId int) {
			defer wg.Done()
			idx := c.NewThreadRecord(threadId)
			for loop := 1; loop <= 5; loop++ {
				requests := []RequestRecord{
					{RequestIndex: loop, Success: true},
				}
				c.AddLoopResult(idx, loop, requests)
			}
		}(i)
	}

	wg.Wait()

	records := c.GetAllRecords()
	if len(records) != numThreads {
		t.Errorf("Expected %d thread records, got %d", numThreads, len(records))
	}

	for _, record := range records {
		if len(record.LoopResults) != 5 {
			t.Errorf("Thread %d: expected 5 loops, got %d", record.ThreadId, len(record.LoopResults))
		}
	}
}

func TestGetAllRecordsReturnsCopy(t *testing.T) {
	c := New()
	c.NewThreadRecord(1)

	records1 := c.GetAllRecords()
	records2 := c.GetAllRecords()

	// Should be different slices
	if &records1[0] == &records2[0] {
		t.Error("GetAllRecords should return a copy, not the same slice")
	}
}
