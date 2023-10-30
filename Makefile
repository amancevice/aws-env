LAMBDA_URL := http://localhost:8080/2015-03-31/functions/function/invocations
VERSION    := latest
ARCH       := $(shell go env GOARCH)

build: bin/aws-env \
	bin/aws-env-darwin-$(ARCH) \
	bin/aws-env-linux-$(ARCH) \
	pkg/aws-env-$(VERSION)-linux-$(ARCH).zip \

clean:
	rm -rf bin opt pkg

test: build
	ARCH=$(ARCH) docker compose up --detach
	ARCH=$(ARCH) docker compose exec lambda curl -s -XPOST -d '{}' $(LAMBDA_URL) &> /dev/null
	ARCH=$(ARCH) docker compose exec lambda curl -s -XPOST -d '{}' $(LAMBDA_URL) &> /dev/null
	ARCH=$(ARCH) docker compose logs
	ARCH=$(ARCH) docker compose down

.PHONY: build clean test

bin/aws-env: bin/aws-env-darwin-$(ARCH)
	cp $< $@

bin/aws-env-darwin-%: go.* **/*.go
	GOOS=darwin GOARCH=$* go build -ldflags="-s -w" -o $@

bin/aws-env-linux-%: go.* **/*.go
	GOOS=linux GOARCH=$* go build -ldflags="-s -w" -o $@
	upx $@

pkg/aws-env-$(VERSION)-%.zip: bin/aws-env-% | opt pkg
	cp $< opt/aws-env
	zip $@ opt

opt pkg:
	mkdir -p $@
