FROM golang:1.21-alpine3.19 AS builder

COPY . /github.com/buhaistrikalo/apex-map/telegram-bot
WORKDIR /github.com/buhaistrikalo/apex-map/telegram-bot

RUN go mod download
RUN go get github.com/joho/godotenv
RUN go build -o ./bin/telegram-bot tg-bot.go

FROM alpine:3.19

WORKDIR /root/

COPY --from=0 /github.com/buhaistrikalo/apex-map/telegram-bot/bin/telegram-bot .
COPY ./images ./images
COPY .env .

EXPOSE 7071

CMD ["./telegram-bot"]