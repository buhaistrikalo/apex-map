FROM golang:1.21-alpine3.19 AS builder

COPY . /github.com/buhaistrikalo/apex-map/server
WORKDIR /github.com/buhaistrikalo/apex-map/server

RUN go mod download
RUN go get github.com/joho/godotenv
RUN go build -o ./bin/server server.go

FROM alpine:3.19

WORKDIR /root/

COPY --from=0 /github.com/buhaistrikalo/apex-map/server/bin/server .
COPY .env .

EXPOSE 7070

CMD ["./server"]