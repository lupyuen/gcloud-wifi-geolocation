// gcloud-wifi-geolocation is a Go web application hosted on Google Cloud Standard App Engine that
// renders realtime temperature sensor data and geolocation on a map.  The sensor data and geolocation
// are pushed via HTTPS from thethings.io Cloud Code Trigger "forward_geolocate" and
// Cloud Code Function "geolocate". The map is rendered using Mapbox GL JS.
// For privacy, users are required to specify the Device ID when viewing the app.  Adapted from
// https://github.com/GoogleCloudPlatform/golang-samples/blob/master/appengine/go11x/helloworld/helloworld.go
// https://github.com/GoogleCloudPlatform/golang-samples/blob/master/appengine_flexible/pubsub/pubsub.go
// https://docs.mapbox.com/mapbox-gl-js/example/3d-buildings/
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
	device := r.URL.Query().Get("device")
	if device == "" {
		// Decode the received JSON.
		msg := &DeviceState{}
		if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
			http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
			return
		}
		device = msg.Device
		if device == "" {
			http.Error(w, "missing device", http.StatusBadRequest)
			return
		}
	}
	enc := json.NewEncoder(w)
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
		fmt.Printf("pull state `%s`: tmp=%f, lat=%f, lng=%f, acc=%f\n", state.Device, state.Tmp, state.Latitude, state.Longitude, state.Accuracy)
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
gcloud app logs tail -s default & ; gcloud app browse ; fg
*/
