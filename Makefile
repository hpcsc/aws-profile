build:
	go build -i -ldflags="-X github.com/hpcsc/aws-profile-utils/commands.version=$(version)" -o bin/aws-profile-utils github.com/hpcsc/aws-profile-utils

run:
	./bin/aws-profile-utils $(args)
