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
	config.Version = sarama.V4_1_1_0
	config.Producer.Compression = sarama.CompressionSnappy

	return sarama.NewAsyncProducer(brokers, config)
}

func PushMessageToQueue(producer sarama.AsyncProducer, topic string, message []byte) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	producer.Input() <- msg
}

func ConnectToConsumerGroup(brokersUrl []string, consumerGroup string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategySticky(),
	}
	return sarama.NewConsumerGroup(brokersUrl, consumerGroup, config)
}
