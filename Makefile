# all our targets are phony (no files to check).
.PHONY: test image push-image image-skip-test image-dev artifacts prune

curdir = $(PWD)
export curdir

test:
	docker build . --target tests

# optional arg: "skip-test=1"
image:
	docker build . -t cardinalby/joytotwitter:release --build-arg SKIP_TEST=${skip-test}

push-image:
	echo "${password}" | docker login -u "${username}" --password-stdin &&\
    docker push cardinalby/joytotwitter:release

image-dev:
	docker build . -f "./Dockerfile.dev" -t cardinalby/joytotwitter:dev

# pass "version" argument and optional "skip-test=1"
artifacts:
	docker build . -f "./Dockerfile.artifacts" -t cardinalby/joytotwitter:artifacts &&\
	docker run -it -v $(curdir)/artifacts:/artifacts/ --env VERSION=$(version) --env SKIP_TEST=${skip-test} cardinalby/joytotwitter:artifacts

prune:
	docker system prune -af