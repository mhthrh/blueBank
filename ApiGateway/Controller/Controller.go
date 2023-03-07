package Controller

import (
	"context"
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
)

type authenticate struct {
	module  string
	user    string
	payload string
}

func init() {
	upgrade = websocket.Upgrader{}
}
func New(t *chan Pool.Connection) {
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

	_ = KafkaBroker.CreateTopic("localhost:9092", customer.UserName)

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
func SocketApis(ctx *gin.Context) {
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
	}
	w, r := ctx.Writer, ctx.Request
	upgrade.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	c, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		ctx.JSON(http.StatusUpgradeRequired, "connote upgrade connection to websocket")
		return
	}
	go newSocketProcess(context.Background(), c, ctx.GetHeader("User_User"))
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
			uId, _ := uuid.NewRandom()
			err := cnn.KafkaWriter.Write(KafkaBroker.Message{
				Topic:    "myTopic",
				Key:      uId.String(),
				Value:    msg,
				MetaData: nil,
			})
			if err != nil {
				errorQueue <- err
			}
			messageToWs <- "successfully"
		case msg := <-messageFromQueue:
			messageToWs <- msg.Value
		case <-errorQueue:
			return
		}
	}
}
