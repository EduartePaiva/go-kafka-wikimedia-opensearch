package main

import (
	"fmt"

	"github.com/eduartepaiva/go-kafka-wikimedia-opensearch/internal/kafka"
	"github.com/eduartepaiva/go-kafka-wikimedia-opensearch/internal/wikimedia"
)

const (
	TOPIC = "wikimedia.recentchange"
)

func main() {

	producer, err := kafka.ConnectToProducer([]string{"localhost:9094"})
	if err != nil {
		panic(err)
	}
	defer producer.Close()
	fmt.Println("connected to kafka broker")

	wikimedia.WikimediaProduceKafka(producer, TOPIC)
}
