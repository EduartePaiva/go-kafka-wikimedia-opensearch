package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/eduartepaiva/go-kafka-wikimedia-opensearch/internal/opensearch"
)

type consumerGroupHandler struct {
	openSearch *opensearch.Client
}

// Cleanup implements [sarama.ConsumerGroupHandler].
func (c *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim implements [sarama.ConsumerGroupHandler].
func (c *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		docId := struct {
			Meta struct {
				Id string `json:"id"`
			} `json:"meta"`
		}{}
		if err := json.
			NewDecoder(bytes.NewReader(msg.Value)).
			Decode(&docId); err != nil {
			return err
		}

		err := c.openSearch.AddToIndex(session.Context(), OPENSEARCH_INDEX, docId.Meta.Id, msg.Value)
		if err != nil {
			fmt.Println(err)
			return err
		}
		// fmt.Printf("received message offset: %d | Topic (%s)\n", msg.Offset, msg.Topic)
		session.MarkMessage(msg, "")
	}
	return nil
}

// Setup implements [sarama.ConsumerGroupHandler].
func (c *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}
