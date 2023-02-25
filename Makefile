LAMBDA_URL := http://localhost:8080/2015-03-31/functions/function/invocations

build: aws-secretsmanager-env.zip

test:
	docker compose up --build --detach
	docker compose exec lambda curl -s -XPOST -d '{}' $(LAMBDA_URL) | jq
	docker compose exec lambda curl -s -XPOST -d '{}' $(LAMBDA_URL) | jq
	docker compose logs
	docker compose down

.PHONY: build test

aws-secretsmanager-env: **/*.go
	docker compose up --build --detach
	docker compose cp lambda:/opt/aws-secretsmanager-env .
	docker compose down

aws-secretsmanager-env.zip: aws-secretsmanager-env
	zip -9qr $@ $<
