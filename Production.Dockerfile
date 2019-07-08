FROM golang:alpine

RUN \
  apk update \
  && apk add --no-cache make wget curl gnupg git

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/golock
COPY . .

RUN make build

CMD ["./bin/golock-linux-amd64"]
