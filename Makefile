IMAGE := sonar-api

sonar-api:
	go build -o bin/sonar-api cmd/sonar-api/sonar-api.go

clean:
	rm -fr bin

build-image:
	docker build -t ${IMAGE} -f docker/Dockerfile .

run-image: build-image
	docker run -v $$(pwd)/example:/data -p 6000:6000 ${IMAGE} /opt/sonar/sonar-api --csv-store /data/pings.csv

push-image:
	docker push ${IMAGE}

.PHONY: test
test:
	go test ./...
