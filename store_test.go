package main

import (
	"sync"
	"testing"
)

func TestConcurrentShortenRequests(t *testing.T) {
	testStore := &URLStore{
		urls: make(map[string]string),
	}

	var wg sync.WaitGroup
	results := make(chan string, 10000)
	errors := make(chan error, 10000)

	// Launch 10,000 concurrent requests
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			code, err := testStore.ShortenRequest("https://github.com")
			if err != nil {
				errors <- err
				return
			}
			results <- code
		}()
	}

	wg.Wait()
	close(results)
	close(errors)

	// Count results
	successCount := len(results)
	errorCount := len(errors)

	if errorCount > 0 {
		t.Fatalf("Got %d errors out of 10000 requests", errorCount)
	}

	// Verify all codes are unique
	seen := make(map[string]bool)
	for code := range results {
		if seen[code] {
			t.Fatalf("Duplicate code generated: %s", code)
		}
		seen[code] = true
	}

	t.Logf("Successfully generated %d unique codes. Store size: %d", successCount, testStore.StoreSize())
}
