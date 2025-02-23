FROM golang:1.24 AS builder

WORKDIR /app

RUN go install github.com/swaggo/swag/v2/cmd/swag@v2.0.0-rc4
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init -g cmd/main.go --output ./docs --parseDependency --parseInternal

RUN CGO_ENABLED=0 GOOS=linux go build -o sylcot cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/sylcot .
COPY --from=builder /app/docs ./docs

EXPOSE ${API_PORT}

CMD ["./sylcot"]