FROM 93s63uis.mirror.aliyuncs.com/library/centos:7.8.2003 as EC-Tools-Centos7-Base

# Install golang
RUN yum install -y wget
RUN wget -q https://golang.google.cn/dl/go1.17.8.linux-amd64.tar.gz
RUN tar -xzf go1.17.8.linux-amd64.tar.gz -C /usr/local/
RUN cp /usr/local/go/bin/* /usr/bin/
RUN go version

# Install required packages
RUN yum install -y git
RUN yum install -y libvirt-devel
RUN yum install -y gcc
RUN yum install -y rpm-build rpmdevtools


# Build project
FROM EC-Tools-Centos7-Base as EC-Tools-Centos7-Builder

# In order not to use caching
ARG DATE

RUN echo ${DATE}
RUN go env -w GO111MODULE="on" \
    && go env -w GOPROXY="https://mirrors.aliyun.com/goproxy/,direct" \
    && cd /root/ec-tools \
    && sh scripts/build.sh
RUN cd /root/ec-tools \
    && sh scripts/build.sh --rpm
