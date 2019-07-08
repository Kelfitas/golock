FROM golang:stretch

RUN apt-get -y update && \
    apt-get -y install build-essential wget apt-transport-https \
    curl software-properties-common

RUN curl -sL https://deb.nodesource.com/setup_12.x | bash -
RUN apt-get install -y yarn nodejs


RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/tebeka/go2xunit
RUN go get -u github.com/go-delve/delve/cmd/dlv

RUN npm i -g nodemon

WORKDIR /go/src/golock
COPY . .

CMD ["/usr/bin/make", "watch"]
