install:
	go get -v -t -d ./...

build:
	go build -i -ldflags="-X github.com/hpcsc/aws-profile/handlers.version=$(version)" -o bin/aws-profile github.com/hpcsc/aws-profile

test:
	go test -v ./...

run:
	./bin/aws-profile $(args)

upload_to_bintray:
	curl -T ./bin/aws-profile \
		 -uhpcsc:${BINTRAY_API_KEY} \
		 -H "X-Bintray-Package:master" \
		 -H "X-Bintray-Version:$(version)" \
		 https://api.bintray.com/content/hpcsc/aws-profile/aws-profile-$(os)-$(version)

	curl -X POST \
		 -uhpcsc:${BINTRAY_API_KEY} \
		 https://api.bintray.com/content/hpcsc/aws-profile/master/$(version)/publish

all: build run
