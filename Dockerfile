FROM golang:1.12 as build

WORKDIR /tmp/customerService_Core

ENV CGO_ENABLED=0
ENV GOPROXY=https://go.likeli.top

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM maven.jsjit.cn:9911/alpine-cert:1.0 as certs

FROM scratch as prod

COPY --from=certs /etc/localtime /etc/localtime
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /tmp/customerService_Core/app .
COPY --from=build /tmp/customerService_Core/conf.yaml .

EXPOSE 5000/tcp

CMD ["/app"]