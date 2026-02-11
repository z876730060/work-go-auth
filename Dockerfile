FROM golang:1.20-alpine as builder

WORKDIR /app
COPY . .

ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct
RUN go build -o auth main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/auth .

EXPOSE 8091

VOLUME [ "/data":"/root/logs" ]

ENTRYPOINT [ "./auth" ]
