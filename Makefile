LAMBDA_URL := http://localhost:8080/2015-03-31/functions/function/invocations
VERSION    := latest
OS         := darwin
ARCH       := $(shell arch)

build: aws-env \
	pkg/aws-env-$(VERSION)-lambda-$(ARCH).zip \
	pkg/aws-env-$(VERSION)-$(OS)-$(ARCH).tar.gz

clean:
	rm -rf aws-env pkg

test: build
	docker compose up --build --detach
	docker compose exec lambda curl -s -XPOST -d '{}' $(LAMBDA_URL) &> /dev/null
	docker compose exec lambda curl -s -XPOST -d '{}' $(LAMBDA_URL) &> /dev/null
	docker compose logs
	docker compose down

.PHONY: build clean test

aws-env: go.* **/*.go
	go build

pkg/aws-env-$(VERSION)-lambda-$(ARCH).zip: | pkg
	docker compose up --build --detach
	docker compose cp lambda:/tmp/package.zip $@
	docker compose down

pkg/aws-env-$(VERSION)-$(OS)-$(ARCH).tar.gz: aws-env | pkg
	tar czf $@ $<

pkg:
	mkdir -p $@
