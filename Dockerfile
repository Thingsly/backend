# syntax=docker/dockerfile:1
FROM golang:alpine AS builder
WORKDIR $GOPATH/src/app
ADD . ./
ENV GO111MODULE=on
RUN go build -o Mitras-Go .

FROM alpine:latest
LABEL description="Mitras Go Backend"
WORKDIR /go/src/app
RUN apk update && apk add --no-cache tzdata
COPY --from=builder /go/src/app .
EXPOSE 9999
RUN chmod +x Mitras-Go
RUN pwd
RUN ls -lrt
ENTRYPOINT [ "./Mitras-Go" ]