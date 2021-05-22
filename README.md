# server-go-gorm

post - http://localhost:8080/videos

{
"title":"video vitle",
"description" : "video description",vrgr
"url":"http://url.test",
"author": {
"first_name": "TestName",
"last_name": "TestLastName",
"age":123,
"email": "test.mail@gmail.com"
} }


Docker run
-
- docker run --rm -d --name dev-postgres  -e POSTGRES_PASSWORD=postgres -p 5432:5432  arm64v8/postgres
- docker build -t my-golang-app .
- docker images
- docker volume create web-golang-volume
- docker run --rm -d --name my-web-golang-app -p 8080:8080 -v web:/app/dockerVolume my-golang-app




