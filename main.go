package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sync"
)

// DeviceState contains the state of the device: temperature and location.
type DeviceState struct {
	tmp       float64
	latitude  float64
	longitude float64
}

var (
	// deviceStates is the threadsafe map that maps device ID to device state.  It needs to be threadsafe because we will reading and writing concurrently.
	deviceStates sync.Map
)

func main() {
	http.HandleFunc("/", indexHandler)
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
	state := DeviceState{math.NaN(), math.NaN(), math.NaN()}
	state.tmp = 28.1

	// Fetch an item that doesn't exist yet.
	result, ok := deviceStates.Load(deviceID)
	if ok {
		state := result.(*DeviceState)
		fmt.Printf("result: `%f` found for key: `hello`\n", state.tmp)
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
		fmt.Printf("result: `%f` found for key: `hello`\n", state.tmp)
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
