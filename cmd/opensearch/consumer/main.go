package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/eduartepaiva/go-kafka-wikimedia-opensearch/internal/kafka"
	"github.com/eduartepaiva/go-kafka-wikimedia-opensearch/internal/opensearch"
)

const (
	TOPIC             = "wikimedia.recentchange"
	CONSUMER_GROUP_ID = "consumer-opensearch-demo"
)

type consumerGroupHandler struct {
}

// Cleanup implements [sarama.ConsumerGroupHandler].
func (c consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim implements [sarama.ConsumerGroupHandler].
func (c consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("received message offset: %d | Topic (%s) | Message %s \n", msg.Offset, msg.Topic, msg.Value[0:10])
		session.MarkMessage(msg, "")
	}
	return nil
}

// Setup implements [sarama.ConsumerGroupHandler].
func (c consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

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

	group, err := kafka.ConnectToConsumerGroup([]string{"localhost:9094"}, CONSUMER_GROUP_ID)
	if err != nil {
		log.Fatal("error while connecting consumer group: ", err)
	}
	defer group.Close()

	handler := &consumerGroupHandler{}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

Loop:
	for {
		if err := group.Consume(ctx, []string{TOPIC}, handler); err != nil {
			log.Println(err)
		}

		if ctx.Err() != nil {
			fmt.Println("closing consumer group")
			break Loop
		}
	}
	fmt.Println("finished processing")
}
