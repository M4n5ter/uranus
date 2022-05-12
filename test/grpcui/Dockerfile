FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /build
ENV GOPATH /build/gopath

RUN go install github.com/fullstorydev/grpcui/cmd/grpcui@latest


FROM alpine

ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/root/bin

COPY --from=builder /build/gopath/bin/grpcui /usr/local/bin/grpcui

CMD ["grpcui", "-version"]
