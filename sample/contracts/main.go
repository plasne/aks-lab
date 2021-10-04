package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type contract struct {
	Artist  string  `json:"artist"`
	Payment float64 `json:"payment"`
}

var contracts = []contract{
	{"Drake", 0.2},
	{"Taylor Swift", 0.25},
	{"Khalid & Normani", 0.1},
}

func getContractForArtist(w http.ResponseWriter, r *http.Request) {
	// see if the artist has a contract
	artist := r.URL.Query().Get("artist")
	var found *contract
	for _, contract := range contracts {
		if contract.Artist == artist {
			found = &contract
			break
		}
	}

	// create a standard contract if necessary
	if found == nil {
		found = &contract{artist, 0.05}
	}

	// write JSON output
	log.Printf("artist \"%v\" is paid \"%v\".\n", artist, found.Payment)
	bytes, err := json.Marshal(found)
	if err != nil {
		http.Error(w, "the song could not be marshalled.", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, "the song could not be written.", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getContractForArtist(w, r)
		default:
			http.Error(w, "the method is not implemented.", http.StatusNotImplemented)
		}
	})
	log.Printf("listening on port %v...\n", 9200)
	err := http.ListenAndServe(":9200", nil)
	log.Fatal(err)
}
