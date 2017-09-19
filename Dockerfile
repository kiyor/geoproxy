FROM golang
ADD . /go/src/github.com/kiyor/geoproxy
RUN cd /go/src/github.com/kiyor/geoproxy && \
	go get && \
	go install github.com/kiyor/geoproxy

EXPOSE 1080
WORKDIR /go/src/github.com/kiyor/geoproxy
VOLUME ["/config"]
ENTRYPOINT ["/go/bin/geoproxy","-l","0.0.0.0:1080","-c","/config"]
