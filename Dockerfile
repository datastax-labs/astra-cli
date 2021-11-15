FROM golang:1.17.3-alpine as BUILD

COPY / /astra
WORKDIR /astra

RUN ./script/build

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /astra/bin/astra /astra
