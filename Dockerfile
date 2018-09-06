FROM ubuntu:16.04

ENV TZ=Asia/Shanghai

COPY app /

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone && apt update && apt install -y ca-certificates

EXPOSE 5000/tcp

CMD ["/app"]