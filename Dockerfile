FROM ubuntu:25.04

RUN apt-get update
RUN apt-get install -y golang \
        git \
        binutils-gold \
        curl \
        g++ \
        gcc \
        gnupg \
        make \
        libasound2-dev

ENV GOOS=linux
ENV GOARCH=arm64
ENV CGO_ENABLED=1
