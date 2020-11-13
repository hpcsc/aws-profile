FROM golang:1.15.5-buster

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.32.2
RUN apt-get update && apt-get install -y expect
