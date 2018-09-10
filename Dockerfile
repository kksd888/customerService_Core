FROM golang:1.10.3 AS build

WORKDIR /go/src/git.jsjit.cn/customerService/customerService_Core

ADD ./godep /usr/local/bin/

ADD . .

RUN CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo -o app .

FROM alpine:latest as certs

RUN apk --update add ca-certificates

FROM scratch AS prod

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=build /go/src/git.jsjit.cn/customerService/customerService_Core/app .

EXPOSE 5000/tcp

CMD ["/app"]