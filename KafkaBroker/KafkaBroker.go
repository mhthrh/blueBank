package KafkaBroker

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaBroker struct {
	Host     string
	ClientId string
	GroupId  string
}
type Message struct {
	Topic    string
	Key      string
	Value    string
	MetaData []byte
}

type Producer struct {
	*kafka.Producer
}

func NewProducer(k KafkaBroker) (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": k.Host,
			"client.id":         k.ClientId,
			"acks":              "all",
		},
	)
	return producer, err
}

func (p *Producer) Produce11(message Message) error {
	var delivery chan kafka.Event
	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &message.Topic, Partition: kafka.PartitionAny},
		Key:            []byte(message.Key),
		Value:          []byte(message.Value)},
		delivery,
	)
	if err != nil {
		return err
	}
	result := <-delivery
	report := result.(*kafka.Message)
	return report.TopicPartition.Error

}

func NewConsumer(k KafkaBroker) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers": k.Host,
			"group.id":          k.GroupId,
			"auto.offset.reset": "smallest",
		},
	)
	return consumer, err
}
