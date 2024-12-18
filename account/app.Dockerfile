FROM golang:1.23.4-alpine3.21 AS build
LABEL authors="mstare"
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/Mostbesep/microservice-com-temp
COPY go.mod go.sum ./
COPY vendor vendor
COPY account account
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./account/cmd/account

FROM alpine:3.21
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]