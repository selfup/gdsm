FROM golang:alpine AS build

ENV CGO_ENABLED=0 APP=gdsm

COPY . $GOPATH/src/github.com/selfup/$APP

WORKDIR $GOPATH/src/github.com/selfup/$APP

RUN go build cmd/daemon/main.go && mv main /go/bin/$APP

FROM scratch

EXPOSE 8081

COPY --from=build /go/bin/gdsm /go/bin/gdsm

CMD ["/go/bin/gdsm"]
