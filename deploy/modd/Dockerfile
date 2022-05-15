FROM golang:1.18.2-alpine3.15
MAINTAINER m4n5ter

ENV GO111MODULE=on
ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime &&  echo 'Asia/Shanghai' >/etc/timezone
RUN go install github.com/cortesi/modd/cmd/modd@latest

CMD ["modd"]

