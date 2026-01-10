FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apk add --no-cache git

RUN go build -o modernovo-server ./GoBackend/cmd/

CMD ["./modernovo-server"]