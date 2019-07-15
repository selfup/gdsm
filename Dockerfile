FROM golang:alpine AS build

ENV CGO_ENABLED=0 APP=gdsm

COPY . $GOPATH/src/github.com/selfup/$APP

WORKDIR $GOPATH/src/github.com/selfup/$APP

RUN go build -o /go/bin/$APP

FROM scratch

EXPOSE 8080

COPY --from=build /go/bin/gdsm /go/bin/gdsm

CMD ["/go/bin/gdsm"]
