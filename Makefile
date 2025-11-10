docker:
	DOCKER_BUILDKIT=1 docker build --network=host -t lemuriaph/go-owot-log-main:latest .