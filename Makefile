# all our targets are phony (no files to check).
.PHONY: test image image-dev artifacts prune

curdir = $(PWD)
export curdir

test:
	docker build . -t cardinalby/joytotwitter:tests --target tests

image:
	docker build . -t cardinalby/joytotwitter:release

image-dev:
	docker build . -f "./Dockerfile.dev" -t cardinalby/joytotwitter:dev

# pass "version" argument
artifacts:
	docker build . -f "./Dockerfile.artifacts" -t cardinalby/joytotwitter:artifacts &&\
	docker run -it -v $(curdir)/artifacts:/artifacts/ --env VERSION=$(version) cardinalby/joytotwitter:artifacts

prune:
	docker system prune -af