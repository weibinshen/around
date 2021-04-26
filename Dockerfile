FROM golang:1.16
MAINTAINER weibinizer@gmail.com

WORKDIR /go/src/around
ADD . /go/src/around

RUN go get cloud.google.com/go/storage
RUN go get github.com/auth0/go-jwt-middleware
RUN go get github.com/form3tech-oss/jwt-go
RUN go get github.com/gorilla/mux
RUN go get github.com/pborman/uuid
RUN go get github.com/olivere/elastic/v7

EXPOSE 8080
CMD ["/usr/local/go/bin/go", "run", "."]
