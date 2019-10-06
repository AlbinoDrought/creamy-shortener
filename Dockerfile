FROM golang:alpine as builder

RUN apk update && apk add git
COPY . $GOPATH/src/github.com/AlbinoDrought/creamy-shortener
WORKDIR $GOPATH/src/github.com/AlbinoDrought/creamy-shortener

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

RUN go get -d -v && go build -a -installsuffix cgo -o /go/bin/creamy-shortener

FROM scratch

COPY --from=builder /go/bin/creamy-shortener /go/bin/creamy-shortener

ENTRYPOINT ["/go/bin/creamy-shortener"]
