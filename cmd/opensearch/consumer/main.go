package main

import (
	"log"

	"github.com/eduartepaiva/go-kafka-wikimedia-opensearch/internal/opensearch"
)

func main() {
	// create open search client
	oshClient, err := opensearch.NewOpenSearchClient([]string{"http://localhost:9200"})
	if err != nil {
		log.Fatal("error while creating opensearch client: ", err)
	}

	// create index
	err = oshClient.CreateIndex("wikimedia")
	if err != nil {
		log.Fatal("error creating wikimedia index: ", err)
	}

	// create kafka consumer client

	// main code logic

	// close things

}
