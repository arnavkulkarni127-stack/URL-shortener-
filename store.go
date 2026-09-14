package main

import "fmt"

// METHODS for the URLStore struct to read and write to the map of shortened URLs and their corresponding original URLs
func (store *URLStore) GET(code string) (string, bool) { // to read something from the map, we need to lock it first, and then unlock it after we're done reading
	store.mu.Lock()
	defer store.mu.Unlock()
	url, ok := store.urls[code]
	return url, ok
}
func (store *URLStore) SET(code, url string) { // to write something to the map, we need to lock it first, and then unlock it after we're done writing
	store.mu.Lock()
	defer store.mu.Unlock()
	store.urls[code] = url
}

// Method to check if a code exists in the map
func (store *URLStore) ShortenRequest(ogurl string) (string, error) { // to check if a code exists in the map, we need to lock it first, and then unlock it after we're done checking
	store.mu.Lock()
	defer store.mu.Unlock() //  everytime we lock something, we need to unlock it after we're done with it
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		code := generateRandomshortCode(6)          // generate a random short code of length 6
		if _, exists := store.urls[code]; !exists { // check if the code already exists in the map
			store.urls[code] = ogurl // if it doesn't exist, add it to the map
			return code, nil
		}
	}

	return "", fmt.Errorf("failed to generate a unique code after %d attempts", maxRetries)
}
func (store *URLStore) StoreSize() int {
	store.mu.Lock()
	defer store.mu.Unlock()
	return len(store.urls)
}
