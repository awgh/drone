FROM alpine:latest

ENV GOROOT=/usr/lib/go \
    GOPATH=/gopath \
    GOBIN=/gopath/bin \
	  CGO_ENABLED=0 \
    PATH=$PATH:$GOROOT/bin:$GOPATH/bin

RUN mkdir -p /gopath/bin && mkdir -p /gopath/pkg && mkdir -p /gopath/src
	
WORKDIR /gopath/bin
ADD config.json /gopath/bin/config.json

RUN apk add -U git go && \  
  go get -v github.com/awgh/drone/drone && \
  apk del git go && \
  rm -rf /gopath/pkg && \
  rm -rf /gopath/src && \
  rm -rf /var/cache/apk/*

EXPOSE 20001
  
ENTRYPOINT ["/gopath/bin/drone"]

# docker build -t awgh/drone .
