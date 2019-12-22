install:
	go get -v -t -d ./...

build:
	go build -i -ldflags="-X github.com/hpcsc/aws-profile-utils/handlers.version=$(version)" -o bin/aws-profile-utils github.com/hpcsc/aws-profile-utils
	cat bintray-descriptor.json | sed -E 's/AWS_PROFILE_UTILS_VERSION/'${version}'/' | tee bintray-descriptor.json

test:
	go test -v ./...

run:
	./bin/aws-profile-utils $(args)

upload_to_bintray:
	curl -T ./bin/aws-profile-utils \
		 -uhpcsc:${BINTRAY_API_KEY} \
		 -H "X-Bintray-Package:aws-profile-utils-master" \
		 -H "X-Bintray-Version:$(version)" \
		 https://api.bintray.com/content/hpcsc/aws-profile-utils/aws-profile-utils-$(os)-$(version)

	curl -X POST \
		 -uhpcsc:${BINTRAY_API_KEY} \
		 https://api.bintray.com/content/hpcsc/aws-profile-utils/aws-profile-utils-master/$(version)/publish

all: build run
