FROM alpine:3.6

RUN mkdir -p /home/app
WORKDIR /home/app

STOPSIGNAL SIGTERM

RUN addgroup -S appuser
RUN adduser -D -S -s /sbin/nologin -G appuser appuser

COPY ./build/app /home/app/go-migrate
RUN chown -c -R appuser /home/app
RUN chmod u=x,g=x,o=-rwx /home/app/go-migrate

USER appuser

CMD ["./go-migrate", "migration", "sync"]