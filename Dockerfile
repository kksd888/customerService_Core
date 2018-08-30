FROM scratch

COPY app /

EXPOSE 5000/tcp

CMD ["/app"]