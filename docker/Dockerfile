FROM golang

WORKDIR /

RUN apt-get update && apt-get install -y --no-install-recommends imagemagick webp

RUN go get -v github.com/nickalie/go-binwrapper && \
    go get -v github.com/stretchr/testify/assert && \
    go get -v golang.org/x/image/webp

RUN mkdir -p $GOPATH/src/github.com/nickalie/go-imagemagickbin
COPY . $GOPATH/src/github.com/nickalie/go-imagemagickbin
WORKDIR $GOPATH/src/github.com/nickalie/go-imagemagickbin
RUN go test -v ./...