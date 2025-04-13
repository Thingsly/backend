package main

import (
	"strconv"
	"time"
)

func main() {
	go TempHumSensor()
	select {}
}

// Get the message ID
func GetMessageID() string {
	// Get the current Unix timestamp
	timestamp := time.Now().Unix()
	// Convert the timestamp to a string
	timestampStr := strconv.FormatInt(timestamp, 10)
	// Extract the last 7 digits
	messageID := timestampStr[len(timestampStr)-7:]

	return messageID
}
