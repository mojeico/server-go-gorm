package repository

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
	"strconv"
	"time"
)

type VideoLoggingRepository interface {
	LoggingCreate(result string)
	LoggingUpdate(result string)
	LoggingDelete(result string)
	LoggingFindAll(result string)
}

type loggingVideo struct {
	connection *redis.Client `json:"-"`
}

type loggingVideoBody struct {
	Key string

	Operation string
	Result    string
}

func NewVideoLoggingRepository() VideoLoggingRepository {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // host:port of the redis server
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	return loggingVideo{
		connection: client,
	}
}

func (l loggingVideo) LoggingCreate(result string) {

	var m = make(map[string]interface{})
	m["operation"] = "create"
	m["result"] = result

	userDefault := "userDefault"

	key := l.keyHelper(userDefault)

	l.connection.HMSet(key, m)

	l.printLog(key)

}

func (l loggingVideo) LoggingUpdate(result string) {

	var m = make(map[string]interface{})
	m["operation"] = "update"
	m["result"] = result

	userDefault := "userDefault"
	key := l.keyHelper(userDefault)

	l.connection.HMSet(key, m)

	l.printLog(key)
}

func (l loggingVideo) LoggingDelete(result string) {

	var m = make(map[string]interface{})
	m["operation"] = "delete"
	m["result"] = result

	userDefault := "userDefault"

	key := l.keyHelper(userDefault)

	l.connection.HMSet(key, m)

	l.printLog(key)
}

func (l loggingVideo) LoggingFindAll(result string) {

	var m = make(map[string]interface{})
	m["operation"] = "findAll"
	m["result"] = result

	userDefault := "userDefault"
	key := l.keyHelper(userDefault)

	l.connection.HMSet(key, m)

	l.printLog(key)
}

func (l loggingVideo) keyHelper(user string) string {
	data := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	key := user + "-" + data
	return key
}

func (l loggingVideo) printLog(key string) {

	loggingData, _ := l.connection.HGetAll(key).Result()

	loggingBody := new(loggingVideoBody)
	mapstructure.Decode(loggingData, loggingBody)

	body := fmt.Sprintf("Key is %s, operation is %s, result is %s.", key, loggingBody.Operation, loggingBody.Result)

	fmt.Println(body)
}
