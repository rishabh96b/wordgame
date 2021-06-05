package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// WordDataStore holds the words data
type WordDataStore struct {
	wordStore map[string]int
}

// Get the word details from the WordDataStore
func (w *WordDataStore) getWordDetails(word string) (int, error) {
	count, ok := w.wordStore[word]
	if !ok {
		// Initialize the word entry in the datastore
		w.wordStore[word] = 1
		return 0, nil
	}
	w.wordStore[word]++
	return count, nil
}

//Create a new WordDataStore with default values
// func NewWordDataStore() WordDataStore {
// 	return WordDataStore{
// 		wordStore: map[string]int{
// 			"lucky":    1,
// 			"magic":    1,
// 			"word":     1,
// 			"new word": 1,
// 		},
// 	}
// }

//DataStore is an interface for different kind of data source/types
type DataStore interface {
	getWordDetails(word string) (int, error)
}

// Logic struct makes the logic independent of any specific kind
// of data store.
type Controller struct {
	dataStore DataStore
}

//"word" used to hold the count of the words
//Specifically made for response purpose
type Word struct {
	Count int
}

func (c *Controller) getDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("getDetails method of controller callled")
	//Get the word requested by the client
	requestedWord := r.URL.Query().Get("word")
	var wordObj Word
	count, err := c.dataStore.getWordDetails(requestedWord)
	if err != nil {
		log.Fatal("Cannot Process the request", err)
	}
	wordObj.Count = count
	if count <= 0 {
		log.Printf("Word %v is not is data store", requestedWord)
		response, err := json.Marshal(wordObj)
		if err != nil {
			log.Fatal("Something went wrong", err)
		}
		log.Println("Responding client: ", response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(response)
		if err != nil {
			log.Fatal("Unable to write to the response", err)
		}

	} else {
		log.Printf("Word %v has count %v", requestedWord, count)
		response, err := json.Marshal(wordObj)
		if err != nil {
			log.Fatal("Something went wrong", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		log.Println("Responding client: ", response)
		_, err = w.Write(response)
		if err != nil {
			log.Fatal("Unable to write to the response", err)
		}
	}
}

func main() {
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
	http.HandleFunc("/", wordController.getDetails)
	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
