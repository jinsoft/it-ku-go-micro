FROM golang:alpine as builder

ENV GO111MODULE on

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk --no-cache add git

WORKDIR /app/ik-api

COPY . .

RUN go env -w GOPROXY=https://goproxy.io,direct && go mod tidy


RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ik-api


FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update && apk add --no-cache bash ca-certificates

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/ik-api/ik-api .

CMD ["./ik-api"]