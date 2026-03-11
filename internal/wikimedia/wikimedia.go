package wikimedia

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/eduartepaiva/go-kafka-wikimedia-opensearch/internal/kafka"
	"github.com/r3labs/sse/v2"
)

const (
	WIKIMEDIA_URL = "https://stream.wikimedia.org/v2/stream/recentchange"
)

func WikimediaProduceKafka(producer sarama.AsyncProducer, topic string) {

	client := sse.NewClient(WIKIMEDIA_URL, func(c *sse.Client) {
		c.Headers["User-Agent"] = "wikimedia-pet-project/1.0"
	})
	events := make(chan *sse.Event, 100)
	client.SubscribeChan("message", events)
	defer client.Unsubscribe(events)

	timeout := time.After(time.Second * 10)

	fmt.Println("scanning producer")

	go func() {
		for {
			select {
			case err, ok := <-producer.Errors():
				if !ok {
					return
				}
				fmt.Printf("failed to push message: %v\n", err)
			case success, ok := <-producer.Successes():
				if !ok {
					return
				}
				fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, success.Partition, success.Offset)
			}
		}
	}()

Loop:
	for {
		select {
		case <-timeout:
			fmt.Println("Timer ended, finishing loop")
			break Loop
		case msg, ok := <-events:
			if !ok {
				break Loop
			}
			kafka.PushMessageToQueue(producer, topic, msg.Data)
		}
	}
}
