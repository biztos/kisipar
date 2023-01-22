# Dockerfile for Kisipar CentOS Build environment - WORK IN PROGRESS
#
# Assumption: CentOS will be easier to deal with than Ubuntu, which was awful.
# Other ideas: ...?
#
# TO CONNECT:
# docker run -t -i --rm -v ~/go:/go -e PATH='/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin' -u kisipar kisipar /bin/bash

FROM centos
MAINTAINER Kevin Frost https://github.com/biztos/kisipar

ENV GOPATH=/go

RUN yum -y update && yum -y install which openssl openssh git && \
    curl -O http://nginx.org/packages/centos/7/x86_64/RPMS/nginx-1.10.1-1.el7.ngx.x86_64.rpm && \
    rpm -i nginx-1.10.1-1.el7.ngx.x86_64.rpm && \
    curl -O https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz && \
    tar -zxf go1.7.1.linux-amd64.tar.gz && mv go /usr/local && \
    useradd -c 'Kisipar Builder' kisipar
