package local

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LocalCache struct {
	debug bool
	path  string
}

// save to cache
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

// get file from cache
func (c LocalCache) Get(hash string, maxAge int) ([]byte, error) {
	// Create file path
	file := filepath.Join(c.path, hash+".png")

	// Check file existence
	stat, err := os.Stat(file)
	if os.IsNotExist(err) {
		return nil, nil // cache miss without error
	}

	// Check last modification date
	if maxAge > 0 && time.Duration(maxAge)*time.Second < time.Now().Sub(stat.ModTime()) {
		if c.debug {
			log.Printf("Stale file %s - renewing", file)
		}

		err = c.Delete(hash)
		if err != nil {
			return nil, err
		}

		return nil, nil
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

// delete file in cache
func (c LocalCache) Delete(hash string) error {
	// TODO
	return nil
}

// init local cache and populate with data
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
