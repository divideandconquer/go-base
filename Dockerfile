FROM scratch

EXPOSE 8080

ADD /bin/app /app

CMD ["/app"]