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

// return file path for a certain hash value
func (c LocalCache) getPathForHash(hash string) string {
	return filepath.Join(c.path, hash[0:1], hash[1:2], hash[2:3])
}

// save to cache
func (c LocalCache) Save(hash string, data []byte) error {
	// Get file path
	path := c.getPathForHash(hash)

	// check path existence and create if if needed
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, 0775)
	}

	file := filepath.Join(path, hash+".png")

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
	file := filepath.Join(c.getPathForHash(hash), hash+".png")

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

		// ignore errors here
		_ = c.Delete(hash)

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
	// Create file path
	file := filepath.Join(c.getPathForHash(hash), hash+".png")

	// try to delete file
	return os.Remove(file)
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
	if err := os.MkdirAll(fp, 0775); err != nil {
		return nil, err
	}

	cache.path = fp
	cache.debug = debug

	return cache, nil
}
