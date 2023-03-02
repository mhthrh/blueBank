package Controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mhthrh/BlueBank/Entity"
	"github.com/mhthrh/BlueBank/Pool"
	"github.com/mhthrh/BlueBank/Proto/bp.go/ProtoGateway"
	"google.golang.org/grpc/status"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var (
	pool *chan Pool.Connection
)

func init() {
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
	ctx.JSON(http.StatusOK, Entity.Gateway{
		UserName:    result.UserName,
		Password:    result.Password,
		Ips:         result.Ips,
		GatewayName: result.GetGatewayName(),
		Status:      result.Status,
	})
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
