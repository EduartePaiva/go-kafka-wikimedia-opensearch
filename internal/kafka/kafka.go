package kafka

import (
	"github.com/IBM/sarama"
)

func ConnectToProducer(brokers []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Retry.Max = 5

	return sarama.NewAsyncProducer(brokers, config)
}

func PushMessageToQueue(producer sarama.AsyncProducer, topic string, message []byte) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	producer.Input() <- msg
}
