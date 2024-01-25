FROM golang:1.21

ENV ENV_PATH="/api/.env"

WORKDIR /api
COPY . .
RUN go mod download
RUN go mod verify
RUN go mod tidy
RUN go build -C ./cmd/api -o ./../../API
CMD ["./API"]
