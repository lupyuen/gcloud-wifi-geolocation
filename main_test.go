package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusOK,
		)
	}
	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf(
			"unexpected body: got (%v) want (%v)",
			rr.Body.String(),
			expected,
		)
	}
}

func TestPushTmp(t *testing.T) {
	pushJSON(t, `{ "device": "test1", "tmp": 28.2 }`, "\"OK\"")
}

func TestPullTmp(t *testing.T) {
	pullDevice(t, "test1", `{"tmp":28.2}`)
}

func TestPushLocation(t *testing.T) {
	pushJSON(t, `{ "device": "test2", "latitude": 1.2730656999999999, "longitude": 103.8096223, "accuracy": 40.1 }`, "\"OK\"")
}

func TestPushLocationThenTmp(t *testing.T) {
	pushJSON(t, `{ "device": "test2", "tmp": 28.2 }`, "\"OK\"")
}

func TestPushTmpThenLocation(t *testing.T) {
	pushJSON(t, `{ "device": "test1", "latitude": 1.2730656999999999, "longitude": 103.8096223, "accuracy": 40.1 }`, "\"OK\"")
}

func TestPull(t *testing.T) {
	pullDevice(t, "test1", `{"accuracy":40.1,"latitude":1.2730656999999999,"longitude":103.8096223,"tmp":28.2}`)
}

func TestIndexHandlerNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/404", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(indexHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusOK,
		)
	}
}

func pullDevice(t *testing.T, device string, expected string) {
	b := []byte("{}")
	req, err := http.NewRequest("POST", "/pull?device="+device, bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(pullHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status: got (%v) want (%v)", status, http.StatusOK)
	}
	if rr.Body.String() != expected {
		t.Errorf("body: got (%v) want (%v)", rr.Body.String(), expected)
	}
}

func pushJSON(t *testing.T, request string, expected string) {
	b := []byte(request)
	req, err := http.NewRequest("POST", "/push", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(pushHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("status: got (%v) want (%v)", status, http.StatusOK)
	}
	if rr.Body.String() != expected {
		t.Errorf("body: got (%v) want (%v)", rr.Body.String(), expected)
	}
}
