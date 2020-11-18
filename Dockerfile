FROM alpine:3.2
WORKDIR /app
ADD mconfig /app
ENTRYPOINT [ "./mconfig" ]
