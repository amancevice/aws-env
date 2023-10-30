LAMBDA_URL := http://localhost:8080/2015-03-31/functions/function/invocations
VERSION    := latest
GOOS       := $(shell go env GOOS)
GOARCH     := $(shell go env GOARCH)

build: bin/aws-env \
	bin/aws-env-darwin-amd64 \
	bin/aws-env-darwin-arm64 \
	bin/aws-env-linux-amd64 \
	bin/aws-env-linux-arm64 \
	pkg/aws-env-$(VERSION)-linux-arm64.zip \
	pkg/aws-env-$(VERSION)-linux-amd64.zip

clean:
	rm -rf bin pkg

test: build
	docker compose up --detach
	docker compose exec lambda curl -s -XPOST -d '{}' $(LAMBDA_URL) &> /dev/null
	docker compose exec lambda curl -s -XPOST -d '{}' $(LAMBDA_URL) &> /dev/null
	docker compose logs
	docker compose down

.PHONY: build clean test

bin/aws-env: bin/aws-env-darwin-$(GOARCH)
	cp $< $@

bin/aws-env-darwin-%: go.* **/*.go
	GOOS=darwin GOARCH=$* go build -ldflags="-s -w" -o $@

bin/aws-env-linux-%: go.* **/*.go
	GOOS=linux GOARCH=$* go build -ldflags="-s -w" -o $@
	upx $@

pkg/aws-env-$(VERSION)-%.zip: bin/aws-env-% | pkg
	cp $< aws-env
	zip $@ aws-env
	rm aws-env

pkg:
	rm -rf $@
	mkdir -p $@
