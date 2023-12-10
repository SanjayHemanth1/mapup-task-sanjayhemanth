FROM golang:1.21.5
WORKDIR /app
COPY go.mod .
COPY main.go .

RUN go build -o bin .

ENTRYPOINT [ "/app/bin" ]
