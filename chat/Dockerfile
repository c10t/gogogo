FROM golang:1.7.3 AS build

WORKDIR /go/src/

ADD chat/ chat
COPY trace/ trace/

RUN go get -d -v golang.org/x/net/html
RUN go get -d -v github.com/gorilla/websocket
RUN go get -d -v github.com/c10t/gogogo/chat/trace
COPY chat/main.go chat/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app chat/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/app .
COPY --from=0 /go/src/chat/templates/ templates/
CMD ["./app"]
