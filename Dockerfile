FROM golang:1.20-alpine as builder

ENV SRC_DIR=/go/src/github.com/wizzardich/geek-reminder-bot/

WORKDIR $SRC_DIR

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . $SRC_DIR

RUN go build -o /app/linkviewer-reincarnated

FROM golang:1.20-alpine 

LABEL org.opencontainers.image.source="https://github.com/wizzardich/geek-reminder-bot" \
      org.opencontainers.image.title="Geek Reminder Bot" \
      org.opencontainers.image.description="A Telegram Bot backend that serves as a Doodle scheduler" \
      org.opencontainers.image.authors="wizzardich" \
      org.opencontainers.image.licenses="MIT"

COPY --from=builder /app/linkviewer-reincarnated /app/linkviewer-reincarnated
WORKDIR /app

EXPOSE 3000

ENTRYPOINT ["./linkviewer-reincarnated"]