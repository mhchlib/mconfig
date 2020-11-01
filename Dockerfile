FROM alpine:3.2
ADD cmd/mconfig /app
WORKDIR /app
ENTRYPOINT [ "/app/mconfig" ]
