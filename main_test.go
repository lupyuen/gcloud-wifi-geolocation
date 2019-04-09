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
	jsonStr := []byte(`{"device":"test1","tmp":28.2}`)
	//  jsonStr := []byte(`{"device":"test1","latitude":1.2730656999999999,"longitude":103.8096223,"accuracy":40.1}`)
	//  jsonStr := []byte(`{"device":"test1","tmp":28.2,"latitude":1.2730656999999999,"longitude":103.8096223,"accuracy":40.1}`)
	expected := "\"OK\""
	req, err := http.NewRequest("POST", "/push", bytes.NewBuffer(jsonStr))
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

func TestPushLocation(t *testing.T) {
	//  jsonStr := []byte(`{"device":"test1","tmp":28.2}`)
	jsonStr := []byte(`{"device":"test2","latitude":1.2730656999999999,"longitude":103.8096223,"accuracy":40.1}`)
	//  jsonStr := []byte(`{"device":"test1","tmp":28.2,"latitude":1.2730656999999999,"longitude":103.8096223,"accuracy":40.1}`)
	expected := "\"OK\""
	req, err := http.NewRequest("POST", "/push", bytes.NewBuffer(jsonStr))
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
