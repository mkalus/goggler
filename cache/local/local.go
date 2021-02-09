package local

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type LocalCache struct {
	debug bool
	path  string
}

func (c LocalCache) Save(hash string, data []byte) error {
	// Create file path
	file := filepath.Join(c.path, hash+".png")

	if c.debug {
		log.Printf("Writing data to %s", file)
	}

	// try to create local file
	fd, err := os.Create(file)
	if err != nil {
		if c.debug {
			log.Printf("Error creating file %s: %s", file, err)
		}

		return err
	}
	defer fd.Close()

	// write data
	_, err = fd.Write(data)
	if err != nil {
		if c.debug {
			log.Printf("Error writing file %s: %s", file, err)
		}

		return err
	}

	return nil
}

func (c LocalCache) Get(hash string) ([]byte, error) {
	// Create file path
	file := filepath.Join(c.path, hash+".png")

	// Check file existence
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return nil, nil // cache miss without error
	}

	if c.debug {
		log.Printf("Getting data from %s", file)
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		if c.debug {
			log.Printf("Error reading file %s: %s", file, err)
		}

		return nil, err
	}

	return data, nil
}

func InitLocalCache(path string, debug bool) (*LocalCache, error) {
	cache := &LocalCache{}

	// convert to absolute path
	fp, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// check if it exists and it is not a directory
	info, err := os.Stat(fp)
	if os.IsExist(err) && !info.Mode().IsDir() {
		return nil, errors.New("path exists, but is not a directory")
	}

	// try to create path
	if err := os.MkdirAll(filepath.Dir(fp), 0775); err != nil {
		return nil, err
	}

	cache.path = fp
	cache.debug = debug

	return cache, nil
}
