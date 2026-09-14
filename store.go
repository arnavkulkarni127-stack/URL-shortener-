package main

// METHODS for the URLStore struct to read and write to the map of shortened URLs and their corresponding original URLs
func (store *URLStore) GET(code string) (string, bool) { // to read something from the map, we need to lock it first, and then unlock it after we're done reading
	store.mu.RLock()
	defer store.mu.RUnlock()
	url, ok := store.urls[code]
	return url, ok
}
func (store *URLStore) SET(code, url string) { // to write something to the map, we need to lock it first, and then unlock it after we're done writing
	store.mu.Lock()
	defer store.mu.Unlock()
	store.urls[code] = url
}
