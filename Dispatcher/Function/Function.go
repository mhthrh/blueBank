package Function

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mhthrh/BlueBank/Entity"
	"github.com/mhthrh/BlueBank/KafkaBroker"
	"github.com/mhthrh/BlueBank/Pool"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoAccount"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/status"
	"time"
)

var (
	Functions     map[string]func()
	pool          Pool.Connection
	socketMessage Entity.WebsocketMessageRequest
)

func init() {
	Functions = make(map[string]func())
	Functions["Accounts,Create"] = f1
	Functions["Accounts,Balance"] = f2
}

func New(p Pool.Connection, s Entity.WebsocketMessageRequest) {
	pool = p
	socketMessage = s
}

func f1() {
	defer func() {
		_ = pool.KafkaWriter.CloseWriter()
	}()
	gCnn := ProtoAccount.NewAccountServicesClient(pool.GrpcConnection)
	var account Entity.Account
	err := mapstructure.Decode(socketMessage.Payload, &account)
	if err != nil {
		NotFound()
		return
	}
	_, stat := gCnn.Create(context.Background(), &ProtoAccount.CreateRequest{
		UserName: account.CustomerUserName,
	})
	st, ok := status.FromError(stat)
	if !ok {
		val, _ := json.Marshal(Entity.WebsocketMessageResponse{
			Id:       socketMessage.Id,
			DateTime: time.Now(),
			Status:   "rejected",
			Reason:   st.Message(),
		})
		_ = pool.KafkaWriter.Write(KafkaBroker.Message{
			Topic:    socketMessage.UserName,
			Key:      socketMessage.Id.String(),
			Value:    string(val),
			MetaData: nil,
		})
		return
	}
	if st != nil {
		val, _ := json.Marshal(Entity.WebsocketMessageResponse{
			Id:       socketMessage.Id,
			DateTime: time.Now(),
			Status:   "rejected",
			Reason:   st.Message(),
		})
		_ = pool.KafkaWriter.Write(KafkaBroker.Message{
			Topic:    socketMessage.UserName,
			Key:      socketMessage.Id.String(),
			Value:    string(val),
			MetaData: nil,
		})
		return
	}
	val, _ := json.Marshal(Entity.WebsocketMessageResponse{
		Id:       socketMessage.Id,
		DateTime: time.Now(),
		Status:   "accepted",
		Reason:   "create account successfully",
	})
	_ = pool.KafkaWriter.Write(KafkaBroker.Message{
		Topic:    socketMessage.UserName,
		Key:      socketMessage.Id.String(),
		Value:    string(val),
		MetaData: nil,
	})
}
func f2() {
	defer func() {
		_ = pool.KafkaWriter.CloseWriter()
	}()
	gCnn := ProtoAccount.NewAccountServicesClient(pool.GrpcConnection)
	var account Entity.Account
	err := mapstructure.Decode(socketMessage.Payload, &account)
	if err != nil {
		NotFound()
		return
	}
	balance, stat := gCnn.Balance(context.Background(), &ProtoAccount.BalanceRequest{
		UserName:      account.CustomerUserName,
		AccountNumber: account.AccountNumber,
	})
	st, ok := status.FromError(stat)
	if !ok {
		val, _ := json.Marshal(Entity.WebsocketMessageResponse{
			Id:       socketMessage.Id,
			DateTime: time.Now(),
			Status:   "rejected",
			Reason:   st.Message(),
		})
		_ = pool.KafkaWriter.Write(KafkaBroker.Message{
			Topic:    socketMessage.UserName,
			Key:      socketMessage.Id.String(),
			Value:    string(val),
			MetaData: nil,
		})
		return
	}
	if st != nil {
		val, _ := json.Marshal(Entity.WebsocketMessageResponse{
			Id:       socketMessage.Id,
			DateTime: time.Now(),
			Status:   "rejected",
			Reason:   st.Message(),
		})
		_ = pool.KafkaWriter.Write(KafkaBroker.Message{
			Topic:    socketMessage.UserName,
			Key:      socketMessage.Id.String(),
			Value:    string(val),
			MetaData: nil,
		})
		return
	}
	val, _ := json.Marshal(Entity.WebsocketMessageResponse{
		Id:       socketMessage.Id,
		DateTime: time.Now(),
		Status:   "accepted",
		Reason:   fmt.Sprintf("balance of %s is: %s", account.AccountNumber, balance),
	})
	_ = pool.KafkaWriter.Write(KafkaBroker.Message{
		Topic:    socketMessage.UserName,
		Key:      socketMessage.Id.String(),
		Value:    string(val),
		MetaData: nil,
	})
}

func NotFound() {
	defer pool.KafkaWriter.CloseWriter()

	val, _ := json.Marshal(Entity.WebsocketMessageResponse{
		Id:       socketMessage.Id,
		DateTime: time.Now(),
		Status:   "rejected",
		Reason:   "category/method not found",
	})
	_ = pool.KafkaWriter.Write(KafkaBroker.Message{
		Topic:    socketMessage.UserName,
		Key:      socketMessage.Id.String(),
		Value:    string(val),
		MetaData: nil,
	})
}
