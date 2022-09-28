FROM golang:1.16 AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .


# MAC OS M1 APPLE  -----RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o main ./cmd/app/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main ./main.go


RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]


FROM alpine:latest AS production
## We have to copy the output from our
## builder stage to our production stage
COPY --from=builder /app .
## we can then kick off our newly compiled
## binary exectuable!!
CMD ["./main"]
