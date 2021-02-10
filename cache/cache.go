package cache

type Cache interface {
	Save(hash string, data []byte) error
	Get(hash string, maxAge int) ([]byte, error)
	Delete(hash string) error
}

// TODO: cache cleaning, max age, etc.
