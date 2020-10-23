FROM alpine:3.2
ADD auth /app
WORKDIR /app
ENTRYPOINT [ "/app/mconfig" ]
