package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type song struct {
	Id     int    `json:"id"`
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Genre  string `json:"genre"`
}

var songs = []song{
	{0, "Drake", "In My Feelings", "HipHop"},
	{1, "Maroon 5", "Girls Like You", "Pop"},
	{2, "Cardi B", "I Like It", "HipHop"},
	{3, "6ix9ine", "FEFE", "Pop"},
	{4, "Post Malone", "Better Now", "Rap"},
	{5, "Eminem", "Lucky You", "Rap"},
	{6, "Juice WRLD", "Lucid Dreams", "Rap"},
	{7, "Eminem", "The Ringer", "Rap"},
	{8, "Travis Scott", "Sicko Mode", "HipHop"},
	{9, "Tyga", "Taste", "HipHop"},
	{10, "Khalid & Normani", "Love Lies", "HipHop"},
	{11, "5 Seconds Of Summer", "Youngblood", "Pop"},
	{12, "Ella Mai", "Boo'd Up", "HipHop"},
	{13, "Ariana Grande", "God Is A Woman", "Pop"},
	{14, "Imagine Dragons", "Natural", "Rock"},
	{15, "Ed Sheeran", "Perfect", "Pop"},
	{16, "Taylor Swift", "Delicate", "Pop"},
	{17, "Florida Georgia Line", "Simple", "Country"},
	{18, "Luke Bryan", "Sunrise, Sunburn, Sunset", "Country"},
	{19, "Jason Aldean", "Drowns The Whiskey", "Country"},
	{20, "Childish Gambino", "Feels Like Summer", "HipHop"},
	{21, "Weezer", "Africa", "Rock"},
	{22, "Panic! At The Disco", "High Hopes", "Rock"},
	{23, "Eric Church", "Desperate Man", "Country"},
	{24, "Nicki Minaj", "Barbie Dreams", "Rap"},
}

var songMutex sync.RWMutex

func retrieve(w http.ResponseWriter, r *http.Request) {
	// use a mutex to safely read from the songs
	songMutex.RLock()
	defer songMutex.RUnlock()

	// get a valid id
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "a valid ID was not provided.", http.StatusBadRequest)
		return
	}
	if id < 0 || id >= len(songs) {
		http.Error(w, "the ID was out-of-range.", http.StatusBadRequest)
		return
	}

	// write JSON output
	log.Printf("retrieving song id %v.\n", id)
	bytes, err := json.Marshal(songs[id])
	if err != nil {
		http.Error(w, "the song could not be marshalled.", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, "the song could not be written.", http.StatusInternalServerError)
		return
	}
}

func store(w http.ResponseWriter, r *http.Request) {
	// append the song
	var val song
	err := json.NewDecoder(r.Body).Decode(&val)
	if err != nil {
		http.Error(w, "the body could not be decoded.", http.StatusBadRequest)
		return
	}

	// use a mutex to protect a change to the songs
	songMutex.Lock()
	val.Id = len(songs)
	songs = append(songs, val)
	songMutex.Unlock()

	// write JSON output
	log.Printf("storing song id %v.\n", val.Id)
	bytes, err := json.Marshal(val)
	if err != nil {
		http.Error(w, "the song could not be marshalled.", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, "the song could not be returned.", http.StatusInternalServerError)
		return
	}
}

func main() {
	godotenv.Load()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			retrieve(w, r)
		case "POST":
			store(w, r)
		default:
			http.Error(w, "the method is not implemented.", http.StatusNotImplemented)
		}
	})
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 80
	}
	log.Printf("listening on port %v...\n", port)
	err = http.ListenAndServe(fmt.Sprint(":", port), nil)
	log.Fatal(err)
}
