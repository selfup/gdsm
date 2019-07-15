FROM golang:alpine

ENV CGO_ENABLED=0 APP=gdsm

COPY . $GOPATH/src/github.com/selfup/$APP

WORKDIR $GOPATH/src/github.com/selfup/$APP

RUN go build -o /go/bin/$APP

EXPOSE 8081

CMD ["/go/bin/gdsm"]
