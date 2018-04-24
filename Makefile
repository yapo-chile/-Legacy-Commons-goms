include scripts/commands/vars.sh


info:
	@echo "YO           : "${YO}
	@echo "ServerRoot   : "${SERVER_ROOT}
	@echo "API Base URL : "${BASE_URL}

setup:
	@scripts/commands/setup.sh

build:
	@scripts/commands/build.sh

start: build run

docker-build:
	@scripts/commands/docker-build.sh

docker-publish:
	@scripts/commands/docker-publish.sh

docker-attach:
	@scripts/commands/docker-attach.sh

docker-compose-up-prod:
	@scripts/commands/docker-compose-up.sh prod

docker-compose-down-prod:
	@scripts/commands/docker-compose-down.sh prod

docker-compose-up-dev:
	@scripts/commands/docker-compose-up.sh dev
	@scripts/commands/docker-attach.sh

docker-compose-down-dev:
	@scripts/commands/docker-compose-down.sh dev

run:
	@${EXEC}

start: build run

test:
	@scripts/commands/test.sh

fix-format:
	@scripts/commands/fix-format.sh
