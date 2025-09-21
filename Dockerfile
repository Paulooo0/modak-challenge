FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download && go mod verify

COPY . .

RUN apk add --no-cache git && \
    go install github.com/swaggo/swag/cmd/swag@v1.16.4 && \
    /go/bin/swag init -g cmd/server/main.go -o docs

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/server ./cmd/server/main.go

FROM scratch AS runner

WORKDIR /app

COPY --from=builder /bin/server /app/server

EXPOSE 8080

CMD ["/app/server"]
