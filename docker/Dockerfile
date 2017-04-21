FROM golang

WORKDIR /

RUN apt-get update && apt-get install -y --no-install-recommends libwebp-dev libjpeg-dev libpng12-dev wget

RUN wget https://www.imagemagick.org/download/ImageMagick.tar.gz && \
    tar xvzf ImageMagick.tar.gz && \
    rm ImageMagick.tar.gz && \
    cd ImageMagick-7.0.5-4 && \
    ./configure && \
    make -j4 && \
    make install && \
    ldconfig /usr/local/lib && \
    cd .. && rm -rf ImageMagick-7.0.5-4

RUN go get -v github.com/nickalie/go-binwrapper && \
    go get -v github.com/stretchr/testify/assert && \
    go get -v golang.org/x/image/webp

RUN mkdir -p $GOPATH/src/github.com/nickalie/go-imagemagickbin
COPY . $GOPATH/src/github.com/nickalie/go-imagemagickbin
WORKDIR $GOPATH/src/github.com/nickalie/go-imagemagickbin
RUN go test -v ./...