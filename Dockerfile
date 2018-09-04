FROM ubuntu:16.04

COPY app /

RUN apt update && apt install -y ca-certificates

EXPOSE 5000/tcp

CMD ["/app"]