FROM golang:1.6

ADD . /workplace/src/github.com/docker/go-dockercloud
WORKDIR /workplace/src/github.com/docker/go-dockercloud/dockercloud
ENV GOPATH /workplace
RUN go get -v
