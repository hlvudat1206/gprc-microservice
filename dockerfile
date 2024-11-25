FROM golang:1.23.2

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod .
COPY main.go contact/server/server.go

RUN go build  -o bin .

ENTRYPOINT ["/app/bin"]
