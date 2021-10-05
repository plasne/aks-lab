package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type song struct {
	Id      int     `json:"id"`
	Title   string  `json:"title"`
	Artist  string  `json:"artist"`
	Payment float64 `json:"payment"`
	Genre   string  `json:"genre"`
}

type contract struct {
	Artist  string  `json:"artist"`
	Payment float64 `json:"payment"`
}

var songsBaseUrl string = "http://songs"
var contractsBaseUrl string = "http://contracts"

func retrieveSong(w http.ResponseWriter, r *http.Request) {
	// call "song" entity service
	songUrl := fmt.Sprint(songsBaseUrl, "/?id=", r.URL.Query().Get("id"))
	log.Printf("fetching song from entity service (%v)...\n", songUrl)
	songResp, err := http.Get(songUrl)
	if err != nil {
		http.Error(w, "failed to contact song service.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// write the output
	var song song
	err = json.NewDecoder(songResp.Body).Decode(&song)
	if err != nil {
		http.Error(w, "failed to get song from entity service.", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	log.Println("successfully retrieved song.")

	// call "contracts" entity service
	contractUrl := fmt.Sprint(contractsBaseUrl, "/?artist=", url.QueryEscape(song.Artist))
	log.Printf("fetching contract from entity service (%v)...\n", contractUrl)
	contractResp, err := http.Get(contractUrl)
	if err != nil {
		http.Error(w, "failed to contact entity service.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// extract the payment
	var contract contract
	err = json.NewDecoder(contractResp.Body).Decode(&contract)
	if err != nil {
		http.Error(w, "failed to get contract from entity service.", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	log.Println("successfully retrieved contract.")
	song.Payment = contract.Payment

	// write the output
	bytes, err := json.Marshal(song)
	if err != nil {
		http.Error(w, "the song could not be marshalled.", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, "the song could not be returned.", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func storeSong(w http.ResponseWriter, r *http.Request) {
	// call "song" entity service
	log.Println("federating store-song request to entity service...")
	resp, err := http.Post(songsBaseUrl, "application/json", r.Body)
	if err != nil {
		http.Error(w, "failed to contact song service.", http.StatusInternalServerError)
		return
	}

	// write the output
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to get song from song service.", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, "the song could not be returned.", http.StatusInternalServerError)
		return
	}
}

func main() {
	// load variables
	godotenv.Load()
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 9000
	}
	if url, ok := os.LookupEnv("SONGS_BASE_URL"); ok {
		songsBaseUrl = url
	}
	if url, ok := os.LookupEnv("CONTRACTS_BASE_URL"); ok {
		contractsBaseUrl = url
	}

	// setup http handlers
	http.HandleFunc("/song", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			retrieveSong(w, r)
		case "POST":
			storeSong(w, r)
		default:
			http.Error(w, "the method is not implemented.", http.StatusNotImplemented)
		}
	})

	// listen
	log.Printf("listening on port %v...\n", port)
	err = http.ListenAndServe(fmt.Sprint(":", port), nil)
	log.Fatal(err)
}
