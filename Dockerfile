FROM golang:1.24 AS swagger-builder

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN ls -la /app && ls -la /app/cmd && ls -la /app/controllers && ls -la /app/models

RUN swag init -g cmd/main.go --output ./docs --parseDependency --parseInternal

RUN CGO_ENABLED=0 GOOS=linux go build -o sylcot cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=swagger-builder /app/sylcot .
COPY --from=swagger-builder /app/docs ./docs

EXPOSE ${API_PORT}

CMD ["./sylcot"]