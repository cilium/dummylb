FROM alpine:3.7

WORKDIR /
COPY dummylb .

CMD ["/dummylb"]
