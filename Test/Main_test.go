package Test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mhthrh/BlueBank/Config"
	"github.com/mhthrh/BlueBank/Entity"
	"github.com/mhthrh/BlueBank/Utils/RandomUtil"
	"github.com/spf13/viper"
	"log"
	"os"
	"testing"
)

var (
	cnn           *sql.DB
	user          Entity.Customer
	login         Entity.CustomerLogin
	gatewayLogin  Entity.GatewayLogin
	bilan         Entity.Bilan
	account       Entity.Account
	ctx           context.Context
	redisClient   *redis.Client
	username      string
	password      string
	accountNumber int64
)

func init() {
	err := os.Chdir("..")
	if err != nil {
		log.Fatal(err)
	}
	cfg := Config.New("Coded.dat", "json", "./ConfigFiles")
	if err := cfg.Initialize(); err != nil {
		log.Fatalf("unable to fill viber, %v", err)
	}
	username = RandomUtil.RandomString(10)
	password = RandomUtil.RandomString(10)
	accountNumber = RandomUtil.RandomInt(10000, 100000)
	user = Entity.Customer{
		FullName: RandomUtil.RandomString(10),
		UserName: username,
		PassWord: password,
		Email:    fmt.Sprintf("%s@gmail.com", RandomUtil.RandomString(10)),
	}
	login = Entity.CustomerLogin{
		UserName: username,
		PassWord: password,
	}
	gatewayLogin = Entity.GatewayLogin{
		UserName: "company1",
		Password: "kir_Khar_Koskesh",
	}
	bilan = Entity.Bilan{
		Name:     RandomUtil.RandomString(10),
		Number:   RandomUtil.RandomInt(10000, 100000),
		Amount:   RandomUtil.RandomInt(10000, 100000),
		IsCredit: true,
	}
	ctx = context.Background()
}
func TestMain(m *testing.M) {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		viper.GetString("Postgresql.Host"), viper.GetInt("Postgresql.Port"), viper.GetString("Postgresql.User"), viper.GetString("Postgresql.password"), viper.GetString("Postgresql.Dbname"))
	cnn, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("Redis.Host"), viper.GetInt("Redis.Port")),
		Password: viper.GetString("Redis.password"),
		DB:       viper.GetInt("Redis.database"),
	})
	_, err = redisClient.Ping().Result()
	if err != nil {
		_ = cnn.Close()
		panic(err)
	}

	os.Exit(m.Run())
}
