FROM golang:1.10.3 AS build

WORKDIR /go/src/git.jsjit.cn/customerService/customerService_Core

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

CMD ["./app"]

FROM alpine:latest as certs

RUN apk --update add ca-certificates

FROM scratch AS prod

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/scboffspring/blog-multistage-go/blog-multistage-go .

EXPOSE 5000/tcp

CMD ["/app"]