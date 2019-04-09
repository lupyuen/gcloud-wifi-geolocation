package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sync"
)

// DeviceState contains the state of the device: temperature (deg C), location and geolocation accuracy (in metres).
// The fields must start with an uppercase letter because these fields will be exported for JSON deserialisation.
type DeviceState struct {
	Device    string
	Tmp       float64
	Latitude  float64
	Longitude float64
	Accuracy  float64
}

var (
	// deviceStates is the threadsafe map that maps device ID to device state.  It needs to be threadsafe because we will reading and writing concurrently.
	deviceStates sync.Map
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/push", pushHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func testDeviceStates() {
	deviceID := "hello"

	// Empty state with no values.
	state := DeviceState{"", math.NaN(), math.NaN(), math.NaN(), math.NaN()}
	state.Tmp = 28.1

	// Fetch an item that doesn't exist yet.
	result, ok := deviceStates.Load(deviceID)
	if ok {
		state := result.(*DeviceState)
		fmt.Printf("result: `%f` found for key: `hello`\n", state.Tmp)
	} else {
		fmt.Println("value not found for key: `hello`")
	}

	// Store an item in the map.
	deviceStates.Store(deviceID, &state)
	fmt.Println("added value: `world` for key: `hello`")

	// Fetch the item we just stored.
	result, ok = deviceStates.Load(deviceID)
	if ok {
		state := result.(*DeviceState)
		fmt.Printf("result: `%f` found for key: `hello`\n", state.Tmp)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	testDeviceStates()
	if r.URL.Path == "/" {
		fmt.Fprint(w, "Hello, World!")
		return
	}
	http.NotFound(w, r)
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	/*
		if r.URL.Query().Get("token") != token {
			http.Error(w, "Bad token", http.StatusBadRequest)
		}
	*/
	msg := &DeviceState{"", math.NaN(), math.NaN(), math.NaN(), math.NaN()}
	if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
		return
	}
	if msg.Device == "" {
		http.Error(w, "Missing device", http.StatusBadRequest)
		return
	}
	if !math.IsNaN(msg.Tmp) {
		fmt.Printf("pushHandler: device=%s, tmp=%f\n", msg.Device, msg.Tmp)
	} else if !math.IsNaN(msg.Latitude) && !math.IsNaN(msg.Longitude) && !math.IsNaN(msg.Accuracy) {
		fmt.Printf("pushHandler: device=%s, lat=%f, lng=%f, acc=%f\n", msg.Device, msg.Latitude, msg.Longitude, msg.Accuracy)
	} else {
		http.Error(w, "Unknown message", http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, "\"OK\"")
}
