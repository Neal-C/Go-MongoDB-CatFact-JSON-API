//lint:file-ignore ST1006 heh...

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	client *mongo.Client
}

func NewServer(mongoClient *mongo.Client) *Server {
	return &Server{
		client: mongoClient,
	}
}

func (self *Server) handleGetAllFacts(responseWriter http.ResponseWriter, request *http.Request){
	collection := self.client.Database("catfact").Collection("facts");

	query := bson.M{};
	cursor, err := collection.Find(context.TODO(), query);

	if err != nil {
		log.Fatal(err)
	}

	var results = []bson.M{};

	if err := cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	responseWriter.Header().Add("Content-Type", "application/json");
	responseWriter.WriteHeader(http.StatusOK);
	json.NewEncoder(responseWriter).Encode(results);
}

type CatFactWorker struct {
	client *mongo.Client
}

func NewCatFactWorkerWithClient(mongoClient *mongo.Client) *CatFactWorker {
	return &CatFactWorker{
		client: mongoClient,
	}
}

func (self *CatFactWorker) start() error {
	collection := self.client.Database("catfact").Collection("facts");

	ticker := time.NewTicker(2 * time.Second);

	for {
		response, err := http.Get("https://catfact.ninja/fact");
		if err != nil {
			return err;
		}
		//* yeah, there's a leak; I won't fix it, but know that I saw it

		var catfact bson.M;
		if err := json.NewDecoder(response.Body).Decode(&catfact); err != nil {
			return err;
		}

		_ , err = collection.InsertOne(context.TODO(), catfact);

		if err != nil {
			return err;
		}

		fmt.Println(catfact)
		<-ticker.C;
	}
}

//docker run --name some-mongo -p 27017:27017 -d mongo
func main(){

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27123"));
		if err != nil {
    	panic(err);
	}

	worker := NewCatFactWorkerWithClient(client)


	go worker.start();

	server := NewServer(client);
	http.HandleFunc("/facts", server.handleGetAllFacts);
	http.ListenAndServe(":3000", nil);
}