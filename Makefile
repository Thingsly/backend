docker-build:
	docker build -t hantdev1/thingsly-go:latest .

docker-push:
	docker push hantdev1/thingsly-go:latest