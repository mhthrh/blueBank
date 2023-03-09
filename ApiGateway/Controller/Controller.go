package Controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mhthrh/BlueBank/Entity"
	"github.com/mhthrh/BlueBank/KafkaBroker"
	"github.com/mhthrh/BlueBank/Pool"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoGateway"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoUser"
	"github.com/mhthrh/BlueBank/Redis"
	"github.com/mhthrh/BlueBank/Token"
	"github.com/spf13/viper"
	"google.golang.org/grpc/status"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var (
	pool          *chan Pool.Connection
	authenticates []authenticate
	upgrade       websocket.Upgrader
	methods       map[string]bool
)

type authenticate struct {
	module  string
	user    string
	payload string
}

func init() {
	methods = make(map[string]bool)
	upgrade = websocket.Upgrader{}
}
func New(t *chan Pool.Connection) {
	fff := viper.Get("Topics")
	print(fff)
	for _, address := range viper.Get("Topics").([]interface{}) {
		topic := address.(map[string]interface{})["name"].(string)
		methods[topic] = true
	}
	pool = t
}

func GatewaySignIn(ctx *gin.Context) {

	var gatewayLogin Entity.GatewayLogin
	if err := ctx.BindJSON(&gatewayLogin); err != nil {
		status, responseError := errorType(err)
		ctx.JSON(status, responseError)
		return
	}
	cnn := getConnection()
	defer func() {
		cnn.Redis.Close()
		cnn.KafkaWriter.CloseWriter()
	}()
	if cnn == nil {
		ctx.JSON(http.StatusInternalServerError, "cannot fetch connection")
		return
	}
	gCnn := ProtoGateway.NewGatewayServicesClient(cnn.GrpcConnection)
	result, stat := gCnn.GatewayLogin(context.Background(), &ProtoGateway.GatewayLoginRequest{
		UserName: gatewayLogin.UserName,
		Password: gatewayLogin.Password,
	})
	st, ok := status.FromError(stat)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, "cannot call api")
		return
	}
	if st != nil {
		ctx.JSON(http.StatusForbidden, st.Message())
		return
	}
	token, err := Token.NewJwtMaker(viper.GetString("SecretKey"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	duration, err := time.ParseDuration(viper.GetString("JWTDuration"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	payload, err := token.Create(gatewayLogin.UserName, duration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, Entity.Gateway{
		UserName:    result.UserName,
		Password:    result.Password,
		Ips:         result.Ips,
		GatewayName: result.GetGatewayName(),
		Status:      result.Status,
		Token:       payload,
	})
}

func UserSignUp(ctx *gin.Context) {

	token, err := Token.NewJwtMaker(viper.GetString("SecretKey"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	payload, err := token.Verify(ctx.GetHeader("Gateway_token"))
	if err != nil {
		ctx.JSON(http.StatusForbidden, "check gateway token")
		return
	}

	if payload.UserName != ctx.GetHeader("Gateway_User") {
		ctx.JSON(http.StatusForbidden, "gateway user mismatch")
		return
	}
	if time.Now().After(payload.ExpireAt) {
		ctx.JSON(http.StatusForbidden, "token is expired")
		return
	}
	var customer Entity.Customer
	if err := ctx.BindJSON(&customer); err != nil {
		s, responseError := errorType(err)
		ctx.JSON(s, responseError.Error())
		return
	}
	cnn := getConnection()
	defer func() {
		cnn.Redis.Close()
		cnn.KafkaWriter.CloseWriter()
	}()
	if cnn == nil {
		ctx.JSON(http.StatusInternalServerError, "cannot fetch connection from pool")
		return
	}

	gCnn := ProtoUser.NewServicesClient(cnn.GrpcConnection)

	_, stat := gCnn.CreateUser(context.Background(), &ProtoUser.UserRequest{
		FullName: customer.FullName,
		UserName: customer.UserName,
		PassWord: customer.PassWord,
		Email:    customer.Email,
	})
	st, ok := status.FromError(stat)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, "cannot call api")
		return
	}
	if st != nil {
		ctx.JSON(http.StatusBadRequest, st.Message())
		return
	}

	_ = KafkaBroker.CreateTopic("localhost:9092", customer.UserName, 1)

	ctx.JSON(http.StatusOK, "create customer successfully")

}
func UserSignIn(ctx *gin.Context) {
	token, err := Token.NewJwtMaker(viper.GetString("SecretKey"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	payload, err := token.Verify(ctx.GetHeader("Gateway_token"))
	if err != nil {
		ctx.JSON(http.StatusForbidden, "check gateway token")
		return
	}

	if payload.UserName != ctx.GetHeader("Gateway_User") {
		ctx.JSON(http.StatusForbidden, "gateway user mismatch")
		return
	}
	if time.Now().After(payload.ExpireAt) {
		ctx.JSON(http.StatusForbidden, "token is expired")
		return
	}
	var customer Entity.CustomerLogin
	if err := ctx.BindJSON(&customer); err != nil {
		s, responseError := errorType(err)
		ctx.JSON(s, responseError.Error())
		return
	}
	cnn := getConnection()
	defer func() {
		cnn.Redis.Close()
		cnn.KafkaWriter.CloseWriter()
	}()
	if cnn == nil {
		ctx.JSON(http.StatusInternalServerError, "cannot fetch connection from pool")
		return
	}
	gCnn := ProtoUser.NewServicesClient(cnn.GrpcConnection)
	_, stat := gCnn.LoginUser(context.Background(), &ProtoUser.LoginRequest{
		UserName: customer.UserName,
		PassWord: customer.PassWord,
	})
	st, ok := status.FromError(stat)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, "cannot call api")
		return
	}
	if st != nil {
		ctx.JSON(http.StatusForbidden, st.Message())
		return
	}

	duration, err := time.ParseDuration(viper.GetString("JWTDuration"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	newPayload, err := token.Create(customer.UserName, duration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, Entity.CustomerLoginResponse{
		UserName:  customer.UserName,
		Token:     newPayload,
		ValidTill: time.Now().Add(duration).String(),
	})
}
func Websocket(ctx *gin.Context) {
	authenticates = append(authenticates, authenticate{
		module:  "Gateway",
		user:    ctx.GetHeader("Gateway_User"),
		payload: ctx.GetHeader("Gateway_Token"),
	})
	authenticates = append(authenticates, authenticate{
		module:  "User",
		user:    ctx.GetHeader("User_User"),
		payload: ctx.GetHeader("User_Token"),
	})
	token, err := Token.NewJwtMaker(viper.GetString("SecretKey"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var durations []time.Duration
	for _, aut := range authenticates {

		payload, err := token.Verify(aut.payload)
		if err != nil {
			ctx.JSON(http.StatusForbidden, fmt.Sprintf("check %s token", aut.module))
			return
		}
		if payload.UserName != aut.user {
			ctx.JSON(http.StatusForbidden, fmt.Sprintf("gateway %s mismatch", aut.module))
			return
		}
		if time.Now().After(payload.ExpireAt) {
			ctx.JSON(http.StatusForbidden, "token is expired")
			return
		}

		durations = append(durations, payload.ExpireAt.Sub(time.Now()))
	}
	w, r := ctx.Writer, ctx.Request
	upgrade.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	wsConnection, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		ctx.JSON(http.StatusUpgradeRequired, "connote upgrade connection to websocket")
		return
	}
	contextUser, _ := context.WithTimeout(context.Background(), durations[0])
	contextGateway, _ := context.WithTimeout(contextUser, durations[1])

	go newSocketProcess(contextGateway, wsConnection, ctx.GetHeader("User_User"))
}

func errorType(e error) (int, error) {
	var sb strings.Builder
	switch reflect.TypeOf(e).String() {
	case "*errors.errorString":
		return http.StatusBadRequest, fmt.Errorf("check message body")
	case "validator.ValidationErrors":
		for _, t := range e.(validator.ValidationErrors) {
			sb.WriteString(fmt.Sprintf("the tag %s validation failed value is: %s", t.StructField(), t.Value()))
		}
		return http.StatusBadRequest, fmt.Errorf(sb.String())
	default:
		return http.StatusBadRequest, fmt.Errorf("cannot deserialize message")
	}
}

func Version(context *gin.Context) {
	context.JSON(http.StatusOK, "Ver:1.0.0")
}
func NotFound(context *gin.Context) {
	context.JSON(http.StatusOK, struct {
		Time        time.Time `json:"time"`
		Description string    `json:"description"`
	}{
		Time:        time.Now(),
		Description: "Workers are working, coming soon!!!",
	})
}
func getConnection() *Pool.Connection {
	select {
	case connection := <-*pool:
		fmt.Println("get new connection")
		return &connection
	case <-time.Tick(time.Second * 1):
		fmt.Println("connection refused")

		return nil
	}
}
func newSocketProcess(ctx context.Context, c *websocket.Conn, userName string) {
	cnn := getConnection()
	reader := KafkaBroker.NewReader([]string{"localhost:9092"}, userName, "groupId-1")
	cntx, cancel := context.WithCancel(context.Background())
	messageFromQueue := make(chan KafkaBroker.Message)
	messageToWs := make(chan string)
	messageFromWs := make(chan string)
	errorQueue := make(chan error, 1)
	defer func() {
		_ = cnn.Redis.Close()
		_ = cnn.KafkaWriter.CloseWriter()
		_ = reader.CloseReader()
		_ = c.Close()
		cancel()
	}()

	go reader.Read(cntx, &messageFromQueue, &errorQueue)
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				errorQueue <- fmt.Errorf("read: %w", err)
				return
			}
			messageFromWs <- string(message)
		}
	}()
	go func() {
		for {
			select {
			case msg := <-messageToWs:
				//string response type is 1
				err := c.WriteMessage(1, []byte(msg))
				if err != nil {
					errorQueue <- fmt.Errorf("write: %w", err)
					return
				}
			}
		}

	}()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-messageFromWs:
			var wsMessage Entity.WebsocketMessageRequest

			err := json.Unmarshal([]byte(msg), &wsMessage)
			if err != nil {
				bytes, _ := json.Marshal(Entity.WebsocketMessageResponse{
					Id:       uuid.UUID{},
					DateTime: time.Now(),
					Status:   "rejected",
					Reason:   err.Error(),
				})
				messageToWs <- string(bytes)
				continue
			}
			wsMessage.UserName = userName
			if !methods[wsMessage.Category] {
				bytes, _ := json.Marshal(Entity.WebsocketMessageResponse{
					Id:       wsMessage.Id,
					DateTime: time.Now(),
					Status:   "rejected",
					Reason:   "method not found",
				})
				messageToWs <- string(bytes)
				continue
			}

			_ = cnn.Redis.Do("SELECT", "1")

			client := Redis.Client{Client: cnn.Redis}
			count, _ := client.KeyExist(wsMessage.Id.String())
			if count > 0 {
				bytes, _ := json.Marshal(Entity.WebsocketMessageResponse{
					Id:       uuid.UUID{},
					DateTime: time.Now(),
					Status:   "rejected",
					Reason:   "duplicate",
				})
				messageToWs <- string(bytes)
				continue
			}
			_ = client.Set(wsMessage.Id.String(), "")
			byt, _ := json.Marshal(&wsMessage)
			err = cnn.KafkaWriter.Write(KafkaBroker.Message{
				Topic:    wsMessage.Category,
				Key:      wsMessage.Id.String(),
				Value:    string(byt),
				MetaData: nil,
			})
			if err != nil {
				bytes, _ := json.Marshal(Entity.WebsocketMessageResponse{
					Id:       wsMessage.Id,
					DateTime: time.Now(),
					Status:   "rejected",
					Reason:   err.Error(),
				})
				messageToWs <- string(bytes)
				continue
			}
			bytes, _ := json.Marshal(Entity.WebsocketMessageResponse{
				Id:       wsMessage.Id,
				DateTime: time.Now(),
				Status:   "received",
				Reason:   "",
			})
			messageToWs <- string(bytes)
		case msg := <-messageFromQueue:
			messageToWs <- msg.Value
		case <-errorQueue:
			return
		}
	}
}
