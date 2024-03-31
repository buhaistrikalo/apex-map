package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"os"

	"github.com/joho/godotenv"
)

func init() {
    // loads values from .env into the system
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
}

func main() {
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
	port := os.Getenv("SERVER_PORT")

	http.HandleFunc("/api/apex-map", apexMapCheckerHandler)
	http.HandleFunc("/api/apex-map/{timestamp}", apexMapCheckerHandler)

	log.Println("Server started at", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func apexMapCheckerHandler(w http.ResponseWriter, r *http.Request) {
	var timestamp int64 
	timeOwn := os.Getenv("SERVER_TIMESTOMP")
	timeOwnInt, err := strconv.ParseInt(timeOwn, 10, 64)
	if err != nil {
		log.Println("Failed to parse timeOwn:", err)
		return
	}

	queryTime := r.URL.Query().Get("time")
	queryTimeInt, err := strconv.ParseInt(queryTime, 10, 64)
	if err != nil {
		log.Println("Failed to parse queryTime:", err)
		// Handle the error accordingly
		return
	}
	if queryTime != "" {
		timestamp = queryTimeInt
	} else {
		timestamp = time.Now().Unix()
	}
    mapDuration := int64(5400) 
    mapNumber := (timestamp - timeOwnInt) / mapDuration

    var currentMap string = getMap(mapNumber)
	var nextMap string = getMap(mapNumber + 1)
	var nextMapTime string = getNextMapTime(mapNumber, mapDuration, timeOwnInt)
	var duration string = getDuration(timestamp, mapNumber, timeOwnInt)

    w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Map    string
		Duration string
		NextMap string
		NextMapTime string
		DateNow string
	}{
		Map:    currentMap,
		Duration: duration,
		NextMap: nextMap,
		NextMapTime: nextMapTime,
		DateNow: time.Now().UTC().Format("2006-01-02 15:04:05"),
	})
}

func getNextMapTime(mapNumber int64, mapDuration int64, timeOwnInt int64) string {
	nextMapTime := time.Unix(timeOwnInt + (mapNumber+1)*mapDuration, 0).UTC().Format("15:04:05")
	return nextMapTime
}

func getDuration (timestamp int64, mapNumber int64, timeOwnInt int64) string {
	eventDuration := int64(5400) 
	var timeRemaining int64 = (mapNumber + 1) * eventDuration - (timestamp - timeOwnInt)
    timeRemainingFormatted := time.Duration(timeRemaining) * time.Second
	return timeRemainingFormatted.String()
}

func getMap(mapNumber int64) string {
	var currentEvent string
	switch mapNumber % 3 {
	case 0:
		currentEvent = "Storm Point"
	case 1:
		currentEvent = "Olympus"
	case 2:
		currentEvent = "Broken Moon"
	}
	return currentEvent
}