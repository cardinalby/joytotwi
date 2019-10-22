# all our targets are phony (no files to check).
.PHONY: test image image-dev artifacts prune

curdir = $(PWD)
export curdir

test:
	go test -v ./app/...

image:
	docker build . -t joytotwitter:release

image-dev:
	docker build . -f "./Dockerfile.dev" -t joytotwitter:dev

# pass "version" argument
artifacts:
	docker build . -f "./Dockerfile.artifacts" -t joytotwitter:artifacts &&\
	docker run -it -v $(curdir)/artifacts:/artifacts/ --env VERSION=$(version) joytotwitter:artifacts

prune:
	docker system prune -af