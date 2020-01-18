run-db:
	docker run --rm -e POSTGRES_PASSWORD=mysecretpassword --network=host postgres

run-api:
	$(shell go env GOPATH)/bin/reflex -d none -s -R vendor. -G ./*/*.go -- go run cmd/api/main.go
