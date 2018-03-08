FROM golang:1.10.0-alpine

RUN apk update && apk add git

RUN mkdir -p /go/src/github.com/byuoitav
ADD . /go/src/github.com/byuoitav/pi-designation-microservice

WORKDIR /go/src/github.com/byuoitav/pi-designation-microservice
RUN go get -d -v
RUN go install -v

CMD ["/go/bin/pi-designation-microservice"]

EXPOSE 5001
