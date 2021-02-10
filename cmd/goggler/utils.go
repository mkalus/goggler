package main

import (
	"github.com/mkalus/goggler/cache"
	"github.com/mkalus/goggler/cache/local"
	"github.com/mkalus/goggler/cache/sthree"
	"github.com/mkalus/goggler/screenshot"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// keeps default settings
var defaultSettings = screenshot.Settings{
	Width:   1920,
	Height:  1024,
	Scale:   0.2,
	Wait:    2000,
	Timeout: 60000,
	Quality: 90,
	MaxAge:  2592000,
}

// listen address for server
var listenAddress = ":8080"

// debug service?
var Debug = false

var MyCache cache.Cache

// collect settings from environment and set them
func defineSettingsFromEnvironment() {
	defaultSettings.Width = getPositiveIntegerFromString(os.Getenv("GOGGLER_WIDTH"), defaultSettings.Width, "GOGGLER_WIDTH", false)
	defaultSettings.Height = getPositiveIntegerFromString(os.Getenv("GOGGLER_HEIGHT"), defaultSettings.Height, "GOGGLER_HEIGHT", false)
	defaultSettings.Scale = getPositiveFloatFromString(os.Getenv("GOGGLER_SCALE"), defaultSettings.Scale, "GOGGLER_SCALE")
	defaultSettings.Wait = getPositiveIntegerFromString(os.Getenv("GOGGLER_WAIT"), defaultSettings.Wait, "GOGGLER_WAIT", false)
	defaultSettings.Timeout = getPositiveIntegerFromString(os.Getenv("GOGGLER_TIMEOUT"), defaultSettings.Timeout, "GOGGLER_TIMEOUT", false)
	defaultSettings.Quality = getPositiveIntegerFromString(os.Getenv("GOGGLER_QUALITY"), defaultSettings.Quality, "GOGGLER_QUALITY", false)
	defaultSettings.MaxAge = getPositiveIntegerFromString(os.Getenv("GOGGLER_MAXAGE"), defaultSettings.MaxAge, "GOGGLER_MAXAGE", true)

	a := os.Getenv("GOGGLER_LISTEN")
	if a != "" {
		listenAddress = a
	}

	d := os.Getenv("GOGGLER_DEBUG")
	if d != "" {
		Debug = true
	}

	if Debug {
		log.Printf("Default settings: width=%dpx, height=%dpx, quality=%d, scale=%f, wait=%dms, timeout=%dms, maxage=%ds",
			defaultSettings.Width,
			defaultSettings.Height,
			defaultSettings.Quality,
			defaultSettings.Scale,
			defaultSettings.Wait,
			defaultSettings.Timeout,
			defaultSettings.MaxAge,
		)
	}

	// init cache
	c := os.Getenv("GOGGLER_CACHE")
	var err error
	switch c {
	case "s3", "S3":
		// implement S3 cache
		url := os.Getenv("GOGGLER_CACHE_S3_URL")
		if url == "" {
			url = "s3.amazonaws.com"
		}
		bucket := os.Getenv("GOGGLER_CACHE_S3_BUCKETNAME")
		accessKey := os.Getenv("GOGGLER_CACHE_S3_ACCESSKEY")
		secretKey := os.Getenv("GOGGLER_CACHE_S3_SECRETKEY")
		region := os.Getenv("GOGGLER_CACHE_S3_REGION")
		ssl := true
		if os.Getenv("GOGGLER_CACHE_S3_SKIPSSL") != "" {
			ssl = false
		}
		createBucket := false
		if os.Getenv("GOGGLER_CACHE_S3_CREATEBUCKET") != "" {
			createBucket = true
		}

		// expiration to days
		days := defaultSettings.MaxAge / 86400

		MyCache, err = sthree.InitS3Cache(url, region, bucket, accessKey, secretKey, days, ssl, createBucket, Debug)
		if err != nil {
			log.Fatal(err)
		}

		// access and secret keys are skipped
		if Debug {
			log.Printf("Cache: s3, url= %s, region=%s, bucket=%s, ssl=%t, expire=%dd", url, region, bucket, ssl, days)
		}
	default:
		// fallback to local cache
		p := os.Getenv("GOGGLER_CACHE_LOCAL_PATH")
		if p == "" {
			p = filepath.Join(os.TempDir(), "goggler")
		}
		MyCache, err = local.InitLocalCache(p, Debug)
		if err != nil {
			log.Fatal(err)
		}

		if Debug {
			log.Printf("Cache: local, path= %s", p)
		}
	}

	// get interval for cleanup runner
	go func() {
		if defaultSettings.MaxAge > 0 {
			i := getPositiveIntegerFromString(os.Getenv("GOGGLER_CACHE_CLEANUP_INTERVAL"), 2592000, "GOGGLER_CACHE_CLEANUP_INTERVAL", true)
			if i > 0 {
				if Debug {
					log.Printf("Cache cleanup: every %s", time.Duration(i)*time.Second)
				}

				for range time.Tick(time.Duration(i) * time.Second) {
					MyCache.RunCleanUp(defaultSettings.MaxAge)
				}
			}
		}
	}()
}

// helper function to parse query string values to positive int
func getPositiveIntegerFromString(n string, defaultValue int, fieldName string, zeroAllowed bool) int {
	// empty value? return default
	if n == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(n)
	if err != nil || i < 0 || (i == 0 && !zeroAllowed) {
		log.Printf("can't convert field %s (value %s) - not a positive integer (falling back to default value)", fieldName, n)
		return defaultValue
	}

	return i
}

// helper function to parse query string values to positive int
func getPositiveFloatFromString(n string, defaultValue float64, fieldName string) float64 {
	// empty value? return default
	if n == "" {
		return defaultValue
	}

	i, err := strconv.ParseFloat(n, 64)
	if err != nil || i <= 0 {
		log.Printf("can't convert field %s (value %s) - not a positive float (falling back to default value)", fieldName, n)
		return defaultValue
	}

	return i
}
