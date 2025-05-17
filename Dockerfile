# syntax=docker/dockerfile:1
FROM --platform=linux/arm64 golang:alpine AS builder
WORKDIR $GOPATH/src/app
ADD . ./
ENV GO111MODULE=on
RUN go build -o thingsly-go .

FROM --platform=linux/arm64 alpine:latest
LABEL description="Thingsly Go Backend"
WORKDIR /go/src/app
RUN apk update && apk add --no-cache tzdata
COPY --from=builder /go/src/app .
EXPOSE 9999
RUN chmod +x thingsly-go
RUN pwd
RUN ls -lrt
ENTRYPOINT [ "./thingsly-go" ]