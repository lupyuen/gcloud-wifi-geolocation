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
	Device    string  `json:"device"` // When we serialise or deserialise to JSON, we will use the "json" names.
	Tmp       float64 `json:"tmp"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
}

var (
	// deviceStates is the threadsafe map that maps device ID to device state.  It needs to be threadsafe because we will reading and writing concurrently.
	deviceStates sync.Map
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/pull", pullHandler)
	http.HandleFunc("/push", pushHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func pullHandler(w http.ResponseWriter, r *http.Request) {
	device := "test1"

	enc := json.NewEncoder(os.Stdout)
	m := map[string]float64{}
	result, ok := deviceStates.Load(device)
	if ok {
		//  State exists for the device. Return the state.
		state := result.(*DeviceState)
		if !math.IsNaN(state.Tmp) {
			m["tmp"] = state.Tmp
		}
		if !math.IsNaN(state.Latitude) {
			m["latitude"] = state.Latitude
		}
		if !math.IsNaN(state.Longitude) {
			m["longitude"] = state.Longitude
		}
		if !math.IsNaN(state.Accuracy) {
			m["accuracy"] = state.Accuracy
		}
	} else {
		fmt.Printf("no state `%s`\n", device)
	}
	enc.Encode(m)
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Verify the token.
	/*
		if r.URL.Query().Get("token") != token {
			http.Error(w, "Bad token", http.StatusBadRequest)
		}
	*/
	// Decode the received JSON.
	msg := &DeviceState{"", math.NaN(), math.NaN(), math.NaN(), math.NaN()}
	if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
		return
	}
	if msg.Device == "" {
		http.Error(w, "Missing device", http.StatusBadRequest)
		return
	}
	// Fetch the current state for the device.
	result, loaded := deviceStates.LoadOrStore(msg.Device, msg)
	if !loaded {
		// State does not exist. We have already stored the message so do nothing.
		fmt.Printf("new state `%s`: tmp=%f, lat=%f, lng=%f, acc=%f\n", msg.Device, msg.Tmp, msg.Latitude, msg.Longitude, msg.Accuracy)
	} else {
		// State already exists. We copy the received message into the current state.
		state := result.(*DeviceState)
		if !math.IsNaN(msg.Tmp) {
			state.Tmp = msg.Tmp
		}
		if !math.IsNaN(msg.Latitude) {
			state.Latitude = msg.Latitude
		}
		if !math.IsNaN(msg.Longitude) {
			state.Longitude = msg.Longitude
		}
		if !math.IsNaN(msg.Accuracy) {
			state.Accuracy = msg.Accuracy
		}
		fmt.Printf("updated state `%s`: tmp=%f, lat=%f, lng=%f, acc=%f\n", state.Device, state.Tmp, state.Latitude, state.Longitude, state.Accuracy)
	}
	fmt.Fprint(w, "\"OK\"")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprint(w, "Hello, World!")
		return
	}
	http.NotFound(w, r)
}

/*
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
*/
