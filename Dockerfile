FROM alpine:3.2
WORKDIR /app
ADD mconfig-server /app
ENTRYPOINT [ "./mconfig-server" ]
