# go commands
build:
	go build -race -o bin/main cmd/app/*.go

run_all:
	go run cmd/app/*.go & go run trucktrace-notification/cmd/app/*.go

run_user_api:
	go run cmd/app/*.go

run_notification_api:
	go run trucktrace-notification/cmd/app/*.go

run_service:
	go run trucktrace-service/cmd/app/*.go
	
clean:
	rm -rf bin/

# docker commands
up:
	sudo docker-compose up -d

down:
	sudo docker-compose down

ps:
	sudo docker ps

exec:
	sudo docker exec -it $(c) bash

screen:
	screen -dmS golang-user go run cmd/app/.*go
	screen -dmS golang-notification go run trucktrace-notification/cmd/app/main.go
	screen -dmS golang-service  go run trucktrace-service/cmd/app/main.go
