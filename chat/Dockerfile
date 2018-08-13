FROM golang:1.7.3 AS build

WORKDIR /go/src/

ADD chat/ chat

RUN go get -d -v golang.org/x/net/html
COPY chat/main.go chat/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app chat/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/app .
CMD ["./app"]