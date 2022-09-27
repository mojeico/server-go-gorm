package queue

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"github.com/trucktrace/pkg/logger"
)

type Message struct {
	Method    string            `json:"method"`
	Body      []byte            `json:"body"`
	Params    map[string]string `json:"params"`
	UserID    uint              `json:"user_id"`
	UserEmail string            `json:"user_email"`
}

type UserNotification struct {
	gorm.Model
	UserKeyData  int    `json:"user_key_data"`
	Message      string `json:"message"`
	Error        bool   `json:"error"`
	ErrorMessage string `json:"error_message"`

	IsRead    bool `json:"is_read"`
	IsDeleted bool `json:"is_deleted"`
}

type SystemNotification struct {
	gorm.Model

	NotificationType string `json:"notification_type"` // expired
	Message          string `json:"message"`
	Error            bool   `json:"error"`
	ErrorMessage     string `json:"error_message"`

	IsDeleted bool `json:"is_deleted"`
}

var natsTime = 10 * time.Minute

type NatsResponse struct {
	Err              error
	UserErrorMessage string
}

func GetNatsConnection() *nats.Conn {

	//host := viper.GetString("nats.host")
	host := os.Getenv("NATS_SERVICE_SERVICE_HOST")
	port := viper.GetString("nats.port")

	natsConnection, err := nats.Connect(fmt.Sprintf("nats://%s:%s", host, port))

	if err != nil {
		panic("GetNatsConnection - Problem with nats connection")
	}

	return natsConnection
}

func RunNatsCreateRequest(connection *nats.Conn, message Message) (NatsResponse, error) {

	byteData, err := json.Marshal(message)

	if err != nil {
		logger.SystemLoggerError("RunNatsCreateRequest", "Can't marshal object ").Error("Error - " + err.Error())
		return NatsResponse{}, err
	}

	msg, err := connection.Request("Create", byteData, natsTime)

	if err != nil {
		logger.SystemLoggerError("RunNatsCreateRequest", "Can't send request to create").Error("Error - " + err.Error())
		return NatsResponse{}, err
	}

	response := NatsResponse{}

	json.Unmarshal(msg.Data, &response)

	return response, response.Err

}

func RunNatsDeleteRequest(connection *nats.Conn, message Message) (NatsResponse, error) {

	byteData, err := json.Marshal(message)

	if err != nil {
		logger.SystemLoggerError("RunNatsDeleteRequest", "Can't marshal object ").Error("Error - " + err.Error())
		return NatsResponse{}, err
	}

	msg, err := connection.Request("Delete", byteData, natsTime)

	if err != nil {
		logger.SystemLoggerError("RunNatsDeleteRequest", "Can't send request to delete").Error("Error - " + err.Error())
		return NatsResponse{}, err
	}

	response := NatsResponse{}

	json.Unmarshal(msg.Data, &response)

	return response, response.Err

}

func RunNatsUpdateRequest(connection *nats.Conn, message Message) (NatsResponse, error) {

	byteData, err := json.Marshal(message)

	if err != nil {
		logger.SystemLoggerError("RunNatsUpdateRequest", "Can't marshal object ").Error("Error - " + err.Error())
		return NatsResponse{}, err
	}

	msg, err := connection.Request("Update", byteData, natsTime)

	if err != nil {
		logger.SystemLoggerError("RunNatsUpdateRequest", "Can't send request to delete").Error("Error - " + err.Error())
		return NatsResponse{}, err
	}

	response := NatsResponse{}

	json.Unmarshal(msg.Data, &response)

	return response, response.Err

}
