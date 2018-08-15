FROM golang:stretch
COPY . $GOPATH/src/freeflix
EXPOSE 8080
ADD https://github.com/ninjaintrouble/freeflix-frontend/releases/download/1.0/frontend.tar $GOPATH/bin
WORKDIR $GOPATH/src/freeflix
RUN apt-get update &&\
    apt-get install gcc &&\
    go install -i -v
WORKDIR $GOPATH/bin
RUN mkdir -p ./torrent/templates &&\
    cp ./../src/freeflix/torrent/templates/status.html ./torrent/templates/status.html
CMD ["freeflix"]