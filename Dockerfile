FROM golang:1.16.4

RUN go get -u github.com/gin-gonic/gin
RUN go get -u github.com/jinzhu/gorm
RUN go get -u github.com/lib/pq


RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o main .



CMD ["/app/main"]