package Test

import (
	"context"
	"fmt"
	"github.com/mhthrh/BlueBank/KafkaBroker"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_KafkaBroker(t *testing.T) {
	message := KafkaBroker.Message{
		Topic:    "myTopic",
		Key:      "one",
		Value:    "kirPolo",
		MetaData: nil,
	}

	writer := KafkaBroker.NewWriter(fmt.Sprintf("%s:%d", viper.GetString("Kafka.Host"), viper.GetInt("Kafka.Port")))
	require.NotNil(t, writer)

	err := writer.Write(message)
	require.NoError(t, err)

	err = writer.CloseWriter()
	require.NoError(t, err)

	address := []string{fmt.Sprintf("%s:%d", viper.GetString("Kafka.Host"), viper.GetInt("Kafka.Port"))}
	topic := "myTopic"
	reader := KafkaBroker.NewReader(address, topic, "testGroupId")
	require.NotNil(t, reader)

	ctxTest, cancel := context.WithCancel(context.Background())
	msg := make(chan KafkaBroker.Message)
	chanError := make(chan error)
	go reader.Read(ctxTest, msg, chanError)

	for range []string{"message", "error"} {
		select {
		case kafkaMessage := <-msg:
			require.Equal(t, kafkaMessage.Key, message.Key)
			require.Equal(t, kafkaMessage.Value, message.Value)
			cancel()
		case e := <-chanError:
			require.Error(t, e)
			err = reader.CloseReader()
			require.NoError(t, err)
		}
	}

}
