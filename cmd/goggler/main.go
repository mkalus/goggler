package main

import (
	"crypto/md5"
	"fmt"
	"github.com/mkalus/goggler/screenshot"
	"log"
	"net/http"
	"net/url"
	"time"
)

func init() {
	defineSettingsFromEnvironment()
}

// start command for goggler
func main() {
	// create a simple web server to handle requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// force base URL only
		if r.URL.Path != "/" {
			w.WriteHeader(404)
			_, _ = w.Write([]byte("404 - not found"))

			if Debug {
				log.Printf("404: %s", r.URL)
			}

			return
		}

		// get query settings
		settings := parseQuery(r.URL)

		// error handling
		if settings.Url == "" {
			w.WriteHeader(500)
			_, _ = w.Write([]byte("500 - query parameter missing, try adding url parameter (other parameters are width, height, scale, quality, wait, and timeout)"))

			if Debug {
				log.Printf("500: %s: query parameter missing", r.URL)
			}

			return
		}

		// TODO: check cache for existing file

		// create screenshot and return it
		image, err := screenshot.CreateScreenShot(settings)
		if err != nil {
			w.WriteHeader(500)
			_, _ = w.Write([]byte(fmt.Sprintf("error occured while creating screenshot: %s", err)))

			if Debug {
				log.Printf("500: %s: %s", r.URL, err)
			}

			return
		}

		// get time passed
		duration := time.Since(start)

		// return image
		log.Printf("200: MISS %s %s", duration, settings.Url)
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(image)
	})

	if Debug {
		log.Printf("Starting server at address %s", listenAddress)
	}

	log.Fatal(http.ListenAndServe(listenAddress, nil))
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
		Width:   getPositiveIntegerFromString(q.Get("width"), defaultSettings.Width, "width"),
		Height:  getPositiveIntegerFromString(q.Get("height"), defaultSettings.Height, "height"),
		Quality: getPositiveIntegerFromString(q.Get("quality"), defaultSettings.Quality, "quality"),
		Scale:   getPositiveFloatFromString(q.Get("scale"), defaultSettings.Scale, "scale"),
		Wait:    getPositiveIntegerFromString(q.Get("wait"), defaultSettings.Wait, "wait"),
		Timeout: getPositiveIntegerFromString(q.Get("timeout"), defaultSettings.Timeout, "timeout"),
	}

	// create file hash
	if settings.Url != "" {
		h := md5.New()
		h.Write([]byte(fmt.Sprintf("%s;%dx%d;%f", settings.Url, settings.Height, settings.Width, settings.Scale)))
		settings.Hash = fmt.Sprintf("%x", h.Sum(nil))
	}

	return settings
}
