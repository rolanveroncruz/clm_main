package main

import (
	"fmt"
	"time"
)

func main() {
	datetimeStr := "2025-05-28 12:31:04 UTC"
	layout := "2006-01-02 15:04:05 MST"

	// Parse the datetime string with the specified layout and UTC location
	parsedTime, err := time.ParseInLocation(layout, datetimeStr, time.UTC)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}

	fmt.Println("Parsed Time (UTC):", parsedTime)

	// Load a different location (e.g., "America/New_York")
	phLocation, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	// Convert the time to the new location
	newYorkTime := parsedTime.In(phLocation)
	fmt.Println("Time in Philippines:", newYorkTime)

	// If you want to keep the time in the original timezone, no conversion is needed:
	originalTime := parsedTime // This is already a time.Time object
	fmt.Println("Original Time (UTC):", originalTime)

}
