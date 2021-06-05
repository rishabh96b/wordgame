package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewWordCountIsZero(t *testing.T) {
	var wordStore DataStore = &WordDataStore{
		wordStore: map[string]int{
			"lucky":    1,
			"magic":    1,
			"word":     1,
			"new word": 1,
		},
	}
	// wordController handles requests based on words
	wordController := Controller{
		dataStore: wordStore,
	}
	request, err := http.NewRequest("GET", "/?word=test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(wordController.getDetails)
	handler.ServeHTTP(rr, request)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	//Checking the response body
	expected := `{"Count":0}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
