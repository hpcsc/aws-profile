install:
	go get gopkg.in/alecthomas/kingpin.v2
	go get gopkg.in/ini.v1

build:
	go build -i -ldflags="-X github.com/hpcsc/aws-profile-utils/handlers.version=$(version)" -o bin/aws-profile-utils github.com/hpcsc/aws-profile-utils

run:
	./bin/aws-profile-utils $(args)

all: build run