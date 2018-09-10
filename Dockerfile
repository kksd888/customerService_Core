FROM golang:1.10.3 AS build

WORKDIR /go/src/git.jsjit.cn/customerService/customerService_Core

ADD ./godep /usr/local/bin/

ADD . .

RUN CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo -o app .

FROM maven.jsjit.cn:9911/ubuntu:16.04.01 as certs

FROM scratch AS prod

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=build /go/src/git.jsjit.cn/customerService/customerService_Core/app .

ENV TZ=Asia/Shanghai

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

EXPOSE 5000/tcp

CMD ["/app"]