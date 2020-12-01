IMAGE=temperature
TAG=$(shell git describe --always --tags)

REGISTRY=registry.cn-shenzhen.aliyuncs.com/xxx

all:
	docker build -t ${REGISTRY}/${IMAGE}:${TAG} .

push: all
	docker push ${REGISTRY}/${IMAGE}:${TAG}

publish: push
	docker tag ${REGISTRY}/${IMAGE}:${TAG} ${REGISTRY}/${IMAGE}:latest
	docker push ${REGISTRY}/${IMAGE}:latest

builder:
	docker build -t go-builder:latest . -f builder.dockerfile
#	docker push ${REGISTRY}/${IMAGE}-builder:latest

test:
	gotest -v ./... || go test -v ./...

lint:
	ls -l | grep '^d' | awk '{print $$NF}' | grep -v vender | xargs golint

count:
	cloc --progress=1 ./ --exclude-dir=vendor,doc,pb

.PHONY: test pb stages



