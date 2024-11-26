FROM golang:1.23.2

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY contact/server/server.go .

RUN go get
RUN go build -o bin .

ENTRYPOINT ["/app/bin"]
