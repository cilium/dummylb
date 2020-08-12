FROM golang:1.14

WORKDIR /dummylb
COPY . .

RUN go build main.go

CMD ["/dummylb/main"]
