package Function

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mhthrh/BlueBank/Entity"
	"github.com/mhthrh/BlueBank/KafkaBroker"
	"github.com/mhthrh/BlueBank/Pool"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoAccount"
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
	_ = json.Unmarshal([]byte(fmt.Sprintf("%v", socketMessage.Payload)), &account)

	_, stat := gCnn.Create(context.Background(), &ProtoAccount.CreateRequest{
		UserName: account.CustomerUserName,
	})
	st, ok := status.FromError(stat)
	if !ok {
		_ = pool.KafkaWriter.Write(KafkaBroker.Message{
			Topic:    socketMessage.UserName,
			Key:      socketMessage.Id.String(),
			Value:    "cannot call api",
			MetaData: nil,
		})
		return
	}
	if st != nil {
		_ = pool.KafkaWriter.Write(KafkaBroker.Message{
			Topic:    socketMessage.UserName,
			Key:      socketMessage.Id.String(),
			Value:    st.Message(),
			MetaData: nil,
		})
		return
	}
}
func f2() {

}

func NotFound() {
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
