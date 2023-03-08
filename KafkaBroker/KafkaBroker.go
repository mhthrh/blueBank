package KafkaBroker

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
	"time"
)

type (
	Writer struct {
		Connection *kafka.Writer
	}
	Reader struct {
		connection *kafka.Reader
	}
	Message struct {
		Topic    string
		Key      string
		Value    string
		MetaData []byte
	}
)

func NewReader(address []string, topic, groupID string) *Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  address,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &Reader{connection: r}

}

func (r *Reader) Read(ctx context.Context, message *chan Message, errChan *chan error) {

	for {
		m, err := r.connection.ReadMessage(ctx)
		if err != nil {
			*errChan <- err
			return
		}
		*message <- Message{
			Topic:    m.Topic,
			Key:      string(m.Key),
			Value:    string(m.Value),
			MetaData: nil,
		}
	}
}
func (r *Reader) CloseReader() error {
	return r.connection.Close()
}
func NewWriter(address string) *Writer {
	return &Writer{
		Connection: &kafka.Writer{
			Addr:     kafka.TCP(address),
			Balancer: &kafka.LeastBytes{},
		},
	}
}
func (w *Writer) Write(msg ...Message) error {
	kafka.TCP()
	var messages []kafka.Message
	for _, message := range msg {
		messages = append(messages, kafka.Message{
			Topic: message.Topic,
			Key:   []byte(message.Key),
			Value: []byte(message.Value),
			Time:  time.Now(),
		})
	}
	err := w.Connection.WriteMessages(context.Background(),
		messages...,
	)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to write messages: %s", err))
	}
	return nil
}
func (w *Writer) CloseWriter() error {
	return w.Connection.Close()
}
func CreateTopic(address, topic string) error {
	conn, err := kafka.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("cannot connect to kafka server, %w", err)
	}
	defer conn.Close()
	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("cannot connect to kafka controller, %w", err)
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return fmt.Errorf("cannot dial with kafka controller, %w", err)
	}
	defer controllerConn.Close()
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return fmt.Errorf("cannot create topic %s with kafka controller, %w", topic, err)
	}
	return nil
}
