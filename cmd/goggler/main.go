package main

import (
	"crypto/md5"
	"fmt"
	"github.com/mkalus/goggler/screenshot"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// start command for goggler
func main() {
	// create a simple web server to handle requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// force base URL only
		if r.URL.Path != "/" {
			w.WriteHeader(404)
			_, _ = w.Write([]byte("404 - not found"))
			return
		}

		// get query settings
		settings := parseQuery(r.URL)

		// error handling
		if settings.Url == "" {
			w.WriteHeader(500)
			_, _ = w.Write([]byte("500 - query parameter missing, try adding url parameter (other parameters are width, height, scale, quality, wait, and timeout)"))
			return
		}

		// TODO: check cache for existing file

		// create screenshot and return it
		image, err := screenshot.CreateScreenShot(settings)
		if err != nil {
			w.WriteHeader(500)
			_, _ = w.Write([]byte(fmt.Sprintf("error occured while creating screenshot: %s", err)))
			return
		}

		// return image
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(image)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// parse query string into settings struct
func parseQuery(r *url.URL) screenshot.Settings {
	// parse values from URL
	if r == nil {
		log.Println("Warning: empty URL")
		return screenshot.Settings{}
	}

	// get query and parse settings
	q := r.Query()

	settings := screenshot.Settings{
		Url:     q.Get("url"),
		Width:   getPositiveIntegerFromString(q.Get("width"), 1920, "width"),
		Height:  getPositiveIntegerFromString(q.Get("height"), 1024, "height"),
		Quality: getPositiveIntegerFromString(q.Get("quality"), 90, "quality"),
		Scale:   getPositiveFloatFromString(q.Get("scale"), 0.1, "scale"),
		Wait:    getPositiveIntegerFromString(q.Get("wait"), 10000, "wait"),
		Timeout: getPositiveIntegerFromString(q.Get("timeout"), 60000, "timeout"),
	}

	// create file hash
	if settings.Url != "" {
		h := md5.New()
		h.Write([]byte(fmt.Sprintf("%s;%dx%d;%f", settings.Url, settings.Height, settings.Width, settings.Scale)))
		settings.Hash = fmt.Sprintf("%x", h.Sum(nil))
	}

	return settings
}

// helper function to parse query string values to positive int
func getPositiveIntegerFromString(n string, defaultValue int, fieldName string) int {
	// empty value? return default
	if n == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(n)
	if err != nil || i <= 0 {
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
