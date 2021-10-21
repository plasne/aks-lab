package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type contract struct {
	Artist  string  `json:"artist"`
	Payment float64 `json:"payment"`
}

var songsBaseUrl string = "http://songs"
var contractsBaseUrl string = "http://contracts"

func retrieveSong(w http.ResponseWriter, r *http.Request) {
	// determine the expected x-api-version
	apiVersion := r.Header.Get("x-api-version")

	// create the request
	songUrl := fmt.Sprint(songsBaseUrl, "/?id=", r.URL.Query().Get("id"))
	songReq, err := http.NewRequest("GET", songUrl, nil)
	if apiVersion != "" {
		songReq.Header.Set("x-api-version", apiVersion)
	}
	if err != nil {
		http.Error(w, "failed to create song request.", http.StatusInternalServerError)
		log.Printf("failed to create song request - %v", err)
		return
	}

	// call "song" entity service
	log.Printf("fetching song from entity service (%v)...\n", songUrl)
	client := http.Client{}
	songResp, err := client.Do(songReq)
	if err != nil {
		http.Error(w, "failed to contact song service.", http.StatusInternalServerError)
		log.Printf("failed to contact song service - %v", err)
		return
	}
	if songResp.StatusCode < 200 || songResp.StatusCode > 299 {
		body, err := io.ReadAll(songResp.Body)
		if err != nil {
			http.Error(w, "received error from song service.", http.StatusInternalServerError)
			log.Printf("received error from song service - %v %v", songResp.StatusCode, err)
		} else {
			http.Error(w, string(body), songResp.StatusCode)
			log.Printf("received error from song service - %v %v", songResp.StatusCode, string(body))
		}
		return
	}

	// write the output
	var song map[string]interface{}
	err = json.NewDecoder(songResp.Body).Decode(&song)
	if err != nil {
		http.Error(w, "failed to get song from entity service.", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	log.Println("successfully retrieved song.")

	// if there is an artist, get the artist's contract
	artistAsInterface, hasArtist := song["artist"]
	artist, artistIsString := artistAsInterface.(string)
	if hasArtist && artistIsString {
		// call "contracts" entity service
		contractUrl := fmt.Sprint(contractsBaseUrl, "/?artist=", url.QueryEscape(artist))
		log.Printf("fetching contract from entity service (%v)...\n", contractUrl)
		contractResp, err := http.Get(contractUrl)
		if err != nil {
			http.Error(w, "failed to contact contracts service.", http.StatusInternalServerError)
			log.Printf("failed to contact contracts service - %v", err)
			return
		}
		if contractResp.StatusCode < 200 || contractResp.StatusCode > 299 {
			body, err := io.ReadAll(contractResp.Body)
			if err != nil {
				http.Error(w, "received error from contracts service.", http.StatusInternalServerError)
				log.Printf("received error from contracts service - %v %v", contractResp.StatusCode, err)
			} else {
				http.Error(w, string(body), contractResp.StatusCode)
				log.Printf("received error from contracts service - %v %v", contractResp.StatusCode, string(body))
			}
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
		song["payment"] = contract.Payment
	}

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
	// determine the expected x-api-version
	apiVersion := r.Header.Get("x-api-version")

	// create the request
	songUrl := fmt.Sprint(songsBaseUrl, "/?id=", r.URL.Query().Get("id"))
	songReq, err := http.NewRequest("POST", songUrl, r.Body)
	songReq.Header.Set("Content-Type", "application/json")
	if apiVersion != "" {
		songReq.Header.Set("x-api-version", apiVersion)
	}
	if err != nil {
		http.Error(w, "failed to create song request.", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// call "song" entity service
	log.Println("federating store-song request to entity service...")
	client := http.Client{}
	resp, err := client.Do(songReq)
	if err != nil {
		http.Error(w, "failed to contact song service.", http.StatusInternalServerError)
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "received error from contracts service.", http.StatusInternalServerError)
			log.Printf("received error from contracts service - %v %v", resp.StatusCode, err)
		} else {
			http.Error(w, string(body), resp.StatusCode)
			log.Printf("received error from contracts service - %v %v", resp.StatusCode, string(body))
		}
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
		port = 80
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
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// returns 200
	})

	// listen
	log.Printf("listening on port %v...\n", port)
	err = http.ListenAndServe(fmt.Sprint(":", port), nil)
	log.Fatal(err)
}
