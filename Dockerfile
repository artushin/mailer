FROM golang:1.7

MAINTAINER Alex Artushin alex@livelyvideo.tv

COPY . /go/src/github.com/LivelyVideo/lively/go

COPY cmd/auth/swagger.json /swagger.json

RUN go install github.com/LivelyVideo/lively/go/cmd/auth

EXPOSE 8080

CMD ["/go/bin/auth"]