# 集成了ca-certificates组件的Ubuntu
FROM maven.jsjit.cn:9911/ubuntu:16.04.01

COPY app /

EXPOSE 5000/tcp

CMD ["/app"]