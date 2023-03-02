package KafkaBroker

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
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

func (r *Reader) Read(ctx context.Context, message chan Message, errChan chan error) {

	for {
		m, err := r.connection.ReadMessage(ctx)
		if err != nil {
			errChan <- err
			return
		}
		message <- Message{
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
