FROM golang:1.15-alpine3.12 as builder

WORKDIR $GOPATH/src/github.com/jicki/alertstore

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && apk upgrade && apk add --no-cache gcc g++ sqlite-libs

ENV GO111MODULE on

ENV GOPROXY https://goproxy.io

COPY . $GOPATH/src/github.com/jicki/alertstore

RUN go mod vendor && go build

# -----------------------------------------------------------------------------

FROM alpine:3.12

LABEL maintainer="jicki"

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && apk upgrade && apk add --no-cache sqlite-libs

WORKDIR /app

COPY --from=builder /go/src/github.com/jicki/alertstore .

EXPOSE 9567

ENTRYPOINT [ "/app/alertstore" ]

