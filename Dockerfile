FROM ubuntu:16.04
MAINTAINER Philip Harries

# RUN echo "deb http://repo.aptly.info/ squeeze main" >> /etc/apt/sources.list.d/aptly.list && \
RUN  echo "deb ftp://ftp.uk.debian.org/debian/ stretch universe main" >> /etc/apt/sources.list.d/golang.list && \
  apt-key adv --keyserver keys.gnupg.net --recv-keys 8B48AD6246925553 && \
  apt-get update && \
  apt-get install -y golang-1.7 git && \
  ln -s /usr/lib/go-1.7/bin/go /usr/bin/go && \
  mkdir -p /root/gowork/src/github.com/PhilipHarries && \
  ln -s  /root/gowork/src/github.com/PhilipHarries/apinate /apinate && \
  GOPATH=/root/gowork go get github.com/mattn/gom

ENV GOPATH /root/gowork
ENV PATH /root/gowork/bin:/usr/bin:/usr/local/bin:/bin:/sbin:/usr/sbin

COPY . /root/gowork/src/github.com/PhilipHarries/apinate

RUN cd /root/gowork/src/github.com/PhilipHarries/apinate && \
  rm -rf /root/gowork/src/github.com/PhilipHarries/apinate/vendor && \
  gom install && \
  gom build -o /bin/apinate

RUN /apinate/run_tests.sh

