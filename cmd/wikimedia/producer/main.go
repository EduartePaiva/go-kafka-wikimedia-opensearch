package main

import (
	"fmt"

	"github.com/eduartepaiva/go-kafka-wikimedia-opensearch/internal"
)

const (
	TOPIC = "wikimedia.recentchange"
)

func main() {

	producer, err := internal.ConnectToProducer([]string{"localhost:9094"})
	if err != nil {
		panic(err)
	}
	defer producer.Close()
	fmt.Println("connected to kafka broker")

	internal.WikimediaProduceKafka(producer, TOPIC)
}
