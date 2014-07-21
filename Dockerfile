FROM ubuntu:trusty

RUN apt-get update

# Install Go
RUN apt-get install -qy build-essential curl git
RUN curl -s https://go.googlecode.com/files/go1.2.1.src.tar.gz | tar -v -C /usr/local -xz
RUN cd /usr/local/go/src && ./make.bash --no-clean 2>&1
RUN mkdir /go
ENV GOPATH /go
ENV PATH /usr/local/go/bin:$PATH
ENV PATH $GOPATH/bin:$PATH

# Install Go dependencies
RUN go get github.com/armon/consul-api
RUN go get github.com/coreos/go-etcd/etcd
RUN go get github.com/fsouza/go-dockerclient

# Install third party tools
RUN apt-get update && apt-get install -y iptables socat

# Compile and install our source
ADD ./start.sh /start.sh
RUN chmod 755 /*.sh

ADD . $GOPATH/src/github.com/jwvdiermen/ambassadord
RUN go install github.com/jwvdiermen/ambassadord

EXPOSE 10000

ENTRYPOINT ["/start.sh"]