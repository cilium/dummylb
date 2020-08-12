FROM golang:1.14

WORKDIR /dummy-lb
COPY . .

RUN go build main.go

CMD ["/dummy-lb/main"]
