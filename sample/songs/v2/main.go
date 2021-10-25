package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type song struct {
	Id     string `json:"id" bson:"_id,omitempty"`
	Artist string `json:"artist" bson:"artist"`
	Title  string `json:"title" bson:"title"`
	Genre  string `json:"genre" bson:"genre"`
}

func retrieve(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	// get a valid id
	id, err := primitive.ObjectIDFromHex(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "a valid ID was not provided.", http.StatusBadRequest)
		log.Printf("a valid ID was not provided - %v", err)
		return
	}

	// get the song from the database
	filter := bson.M{"_id": id}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var val song
	err = collection.FindOne(ctx, filter).Decode(&val)
	if err == mongo.ErrNoDocuments {
		http.Error(w, "no song with that id was found.", http.StatusNotFound)
		log.Printf("the song was not found for id %v.", id)
		return
	} else if err != nil {
		http.Error(w, "the song could not be retrieved.", http.StatusInternalServerError)
		log.Printf("the song could not be retrieved - %v", err)
		return
	}

	// write JSON output
	log.Printf("retrieving song id %v.\n", id)
	bytes, err := json.Marshal(val)
	if err != nil {
		http.Error(w, "the song could not be marshalled.", http.StatusInternalServerError)
		log.Printf("the song could not be marshalled - %v", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, "the song could not be written.", http.StatusInternalServerError)
		log.Printf("the song could not be written - %v", err)
		return
	}
}

func store(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	// decode the input
	var val song
	err := json.NewDecoder(r.Body).Decode(&val)
	if err != nil {
		http.Error(w, "the body could not be decoded.", http.StatusBadRequest)
		return
	}

	// insert into the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, val)
	if err != nil {
		log.Fatalf("failed to add song %v", err)
	}
	val.Id = result.InsertedID.(primitive.ObjectID).Hex()

	// write JSON output
	log.Printf("stored song id %v.\n", val.Id)
	bytes, err := json.Marshal(val)
	if err != nil {
		http.Error(w, "the song could not be marshalled.", http.StatusInternalServerError)
		log.Printf("the song could not be marshalled - %v", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, "the song could not be returned.", http.StatusInternalServerError)
		log.Printf("the song could not be written - %v", err)
		return
	}
}

func EnvOrString(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		val = def
	}
	return val
}

func EnvOrInt(key string, def int) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		val = def
	}
	return val
}

func main() {
	// determine configuration
	godotenv.Load()
	port := EnvOrInt("PORT", 80)
	mongoConnString := EnvOrString("MONGO_CONNSTRING", "")
	if mongoConnString == "" {
		log.Fatal("You must provide MONGO_CONNSTRING.")
	}
	mongoDatabase := EnvOrString("MONGO_DATABASE", "db")
	mongoCollection := EnvOrString("MONGO_COLLECTION", "col")
	log.Printf("PORT = %v", port)
	log.Print("MONGO_CONNSTRING = *SET*")
	log.Printf("MONGO_DATABASE = %v", mongoDatabase)
	log.Printf("MONGO_COLLECTION = %v", mongoCollection)

	// attempt to initialize Cosmos connection
	log.Printf("attempting to initialize Cosmos connection...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnString))
	if err != nil {
		log.Fatalf("unable to initialize Cosmos connection - %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	log.Printf("successfully initialized Cosmos connection.")

	// attempt to connect to a Cosmos instance
	log.Printf("attempting to connect to Cosmos...")
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(pingCtx, nil)
	if err != nil {
		log.Fatalf("unable to connect to Cosmos - %v", err)
	}
	pingCancel()
	log.Println("successfully connected to Cosmos.")
	collection := client.Database(mongoDatabase).Collection(mongoCollection)

	// create HTTP handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			retrieve(w, r, collection)
		case "POST":
			store(w, r, collection)
		default:
			http.Error(w, "the method is not implemented.", http.StatusNotImplemented)
		}
	})

	// start listening for incoming connections
	log.Printf("listening on port %v...", port)
	err = http.ListenAndServe(fmt.Sprint(":", port), nil)
	log.Fatal(err)
}
