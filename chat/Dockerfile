FROM golang:1.7.3 AS build

WORKDIR /go/src/

ADD echo-server/ echo-server

RUN go get -d -v golang.org/x/net/html
COPY echo-server/main.go echo-server/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app echo-server/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/app .
CMD ["./app"]
