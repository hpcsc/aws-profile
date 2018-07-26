build:
	go build -i -ldflags="-X github.com/hpcsc/aws-profile-utils/commands.version=$(version)" -o build/aws-profile-utils github.com/hpcsc/aws-profile-utils

run:
	./aws-profile-utils $(args)
