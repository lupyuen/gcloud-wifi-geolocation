package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type deviceState struct {
}

var (
	//  Threadsafe map that maps device ID to device state.  It needs to be threadsafe because we will reading and writing concurrently.
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
	//  Fetch an item that doesn't exist yet.
	deviceID := "hello"
	result, ok := deviceStates.Load(deviceID)
	if ok {
		fmt.Println(result.(string))
	} else {
		fmt.Println("value not found for key: `hello`")
	}

	//  Store an item in the map.
	deviceStates.Store("hello", "world")
	fmt.Println("added value: `world` for key: `hello`")

	// Fetch the item we just stored.
	result, ok = deviceStates.Load("hello")
	if ok {
		fmt.Printf("result: `%s` found for key: `hello`\n", result.(string))
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
