FROM golang:1.8
MAINTAINER Erick Torres <erick@schibsted.cl>

COPY goms conf/ /home/user/go/src/github.schibsted.io/Yapo/goms
WORKDIR /home/user/go/src/github.schibsted.io/Yapo/goms

CMD ["./goms"]
