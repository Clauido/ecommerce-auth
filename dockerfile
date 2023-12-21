FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN rm -f .env

COPY docker.env .env

RUN go get -d -v ./...

CMD ["go","run","main.go"]
