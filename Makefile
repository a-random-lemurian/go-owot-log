docker:
	DOCKER_BUILDKIT=1 docker build --network=host -t lemuriaph/go-owot-log-main:latest .

run-staging:
	go run ./cmd/owot-log start --config config-staging.yml
