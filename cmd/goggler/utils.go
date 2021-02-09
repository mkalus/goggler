package main

import (
	"github.com/mkalus/goggler/screenshot"
	"log"
	"os"
	"strconv"
)

// keeps default settings
var defaultSettings = screenshot.Settings{
	Width:   1920,
	Height:  1024,
	Scale:   0.2,
	Wait:    2000,
	Timeout: 60000,
	Quality: 90,
}

// listen address for server
var listenAddress = ":8080"

// debug service?
var Debug = false

// collect settings from environment and set them
func defineSettingsFromEnvironment() {
	defaultSettings.Width = getPositiveIntegerFromString(os.Getenv("GOGGLER_WIDTH"), defaultSettings.Width, "GOGGLER_WIDTH")
	defaultSettings.Height = getPositiveIntegerFromString(os.Getenv("GOGGLER_HEIGHT"), defaultSettings.Height, "GOGGLER_HEIGHT")
	defaultSettings.Scale = getPositiveFloatFromString(os.Getenv("GOGGLER_SCALE"), defaultSettings.Scale, "GOGGLER_SCALE")
	defaultSettings.Wait = getPositiveIntegerFromString(os.Getenv("GOGGLER_WAIT"), defaultSettings.Wait, "GOGGLER_WAIT")
	defaultSettings.Timeout = getPositiveIntegerFromString(os.Getenv("GOGGLER_TIMEOUT"), defaultSettings.Timeout, "GOGGLER_TIMEOUT")
	defaultSettings.Quality = getPositiveIntegerFromString(os.Getenv("GOGGLER_QUALITY"), defaultSettings.Quality, "GOGGLER_QUALITY")

	a := os.Getenv("GOGGLER_LISTEN")
	if a != "" {
		listenAddress = a
	}

	d := os.Getenv("GOGGLER_DEBUG")
	if d != "" {
		Debug = true
	}

	if Debug {
		log.Printf("Default settings: width=%dpx, height=%dpx, quality=%d, scale=%f, wait=%dms, timeout=%dms",
			defaultSettings.Width,
			defaultSettings.Height,
			defaultSettings.Quality,
			defaultSettings.Scale,
			defaultSettings.Wait,
			defaultSettings.Timeout,
		)
	}
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
