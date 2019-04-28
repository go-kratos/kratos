FROM golang:1.10.2

RUN mv /etc/apt/sources.list /etc/apt/sources.list.bak \
    && echo deb http://mirrors.aliyun.com/debian stretch main contrib non-free > /etc/apt/sources.list \
    && echo deb-src http://mirrors.aliyun.com/debian stretch main contrib non-free >> /etc/apt/sources.list \
    && echo deb http://mirrors.aliyun.com/debian stretch-updates main contrib non-free >> /etc/apt/sources.list \
    && echo deb-src http://mirrors.aliyun.com/debian stretch-updates main contrib non-free >> /etc/apt/sources.list \
    && echo deb http://mirrors.aliyun.com/debian-security stretch/updates main contrib non-free >> /etc/apt/sources.list \
    && echo deb-src http://mirrors.aliyun.com/debian-security stretch/updates main contrib non-free >> /etc/apt/sources.list \
    && apt-get update

RUN apt-get update && apt-get install openjdk-8-jdk -y

ADD ./ /go/src/go-common/

RUN cd /go/src/go-common/app/tool/kratos && go install