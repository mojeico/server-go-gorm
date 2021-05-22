FROM arm64v8/postgres


RUN mkdir /app
WORKDIR /app

ADD .. /app
RUN go build -o main .



CMD ["/app/main"]