package cache

type Cache interface {
	// Save image data to cache
	Save(hash string, data []byte) error

	// get image data from cache (possibly invalidating and deleting stale data)
	Get(hash string, maxAge int) ([]byte, error)

	// run periodic cleanup service (this can run longer)
	RunCleanUp(maxAge int)
}
