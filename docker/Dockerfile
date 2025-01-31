FROM golang:alpine

RUN apk update && \
    apk add curl

RUN mkdir /app
COPY ./ /app

WORKDIR /app

RUN go install mvdan.cc/gofumpt@latest
RUN go install github.com/bombsimon/wsl/v4/cmd...@master
RUN go install github.com/Antonboom/testifylint@latest
