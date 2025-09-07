FROM golang:1.25-alpine as builder

ENV SRC_DIR=/go/src/github.com/wizzardich/linkviewer-reincarnated/

WORKDIR $SRC_DIR

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . $SRC_DIR

RUN go build -o /app/linkviewer-reincarnated

FROM golang:1.25-alpine 

LABEL org.opencontainers.image.source="https://github.com/wizzardich/linkviewer-reincarnated" \
      org.opencontainers.image.title="Link Viwer Reincarnated" \
      org.opencontainers.image.description="Link viewer backend in Go" \
      org.opencontainers.image.authors="wizzardich" \
      org.opencontainers.image.licenses="MIT"

COPY --from=builder /app/linkviewer-reincarnated /app/linkviewer-reincarnated
WORKDIR /app

EXPOSE 3000

ENTRYPOINT ["./linkviewer-reincarnated"]
