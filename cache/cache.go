package cache

type Cache interface {
	Save(hash string, data []byte) error
	Get(hash string) ([]byte, error)
}

// TODO: cache cleaning, max age, etc.
