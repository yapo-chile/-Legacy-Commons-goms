.PHONY : status build rpm-clean rpm-build
get_pid=$$(ps ux | grep -v grep | grep -E "\./goms" | awk '{ print $$2 }')

EXEC="./goms"
GOLINT=$$GOPATH/bin/golint

GENPORTOFF?=0
genport = $(shell expr ${GENPORTOFF} + \( $(shell id -u) - \( $(shell id -u) / 100 \) \* 100 \) \* 200 + 30100 + $(1))

YO=`whoami`
LISTEN_PORT=$(call genport,2)
SERVER_ROOT=${PWD}
SERVERNAME=`hostname`
BASE_URL="http://"${SERVERNAME}":"${LISTEN_PORT}

info:
	@echo "YO           : "${YO}
	@echo "ServerRoot   : "${SERVER_ROOT}
	@echo "API Base URL : ${BASE_URL}"

build: stop
	go get
	go build

run:
	@PID=$(call get_pid); \
	if [ "$$PID" ]; then \
		echo "ERROR: service is already running (PID: $$PID)"; \
	else if [ ! -x "${EXEC}" ]; then \
		echo "service ${EXEC} not found"; \
	else \
		${MAKE} -s update_config demonize; \
	fi; \
	fi

kill:
	@PID=$(call get_pid); \
	if [ "$$PID" ]; then \
		kill -15 $$PID; \
	fi

start: build run
	@${MAKE} -s status

stop: kill
	@${MAKE} -s status

restart: start

status: service-status

service-status:
	@PID=$(call get_pid); \
	if [ "$$PID" ]; then \
		echo -e "\e[32mSERVICE RUNNING (PID: $$PID)\e[0m"; \
	else \
		echo -e "\e[31mSERVICE NOT RUNNING\e[0m"; \
	fi

update_config:
	@rm -f conf/conf.json
	@cp conf/conf.json.in conf/conf.json
	@sed -i "s/__SERVER_NAME__/${SERVERNAME}/g" conf/conf.json
	@sed -i "s/__SERVICE_PORT__/${LISTEN_PORT}/g" conf/conf.json
	@sed -i "s/__SYSLOG_ENABLED__/false/g" conf/conf.json
	@sed -i "s/__SYSLOG_IDENTITY__/goms/g" conf/conf.json
	@sed -i "s/__STDLOG_ENABLED__/true/g" conf/conf.json
	@sed -i "s/__LOG_LEVEL__/0/g" conf/conf.json

demonize:
	@nohup ${EXEC} >> logs/${EXEC}.log 2>&1 &

rpm-clean:
	@PID=$(call get_pid); \
	if [ "$$PID" ]; then \
		echo "ERROR: service is running (PID: $$PID)"; \
		exit -1; \
	else \
		rm -Rf build; \
	fi

rpm-build: rpm-clean rpm-setuptree build
	cp conf/conf.json.in conf/conf.json
	rpmbuild -bb scripts/api.spec
	mv build/RPMS/x86_64/goms*.rpm ./
	rm -Rf build

rpm-setuptree:
	mkdir -p build/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

check: check-format check-vet check-lint

check-format:
	@echo "==> Checking format with gofmt:"
	@ERRORS=0; \
	for file in $$(find ./ -iname "*.go"); do \
		echo -n "checking $${file:2}" ; \
		errors=$$(gofmt -d -e -s $$file | grep -c -E "^\+[^\+]"); \
		if [ "$$errors" -gt 0 ]; then \
			ERRORS=$$((ERRORS + errors)); \
			echo -e " ... \e[31mfound $$errors issues\e[0m"; \
		else \
			echo " ... ok"; \
		fi; \
	done; \
	if [ $$ERRORS -gt 0 ]; then \
		echo "found $$ERRORS errors"; \
		exit 1; \
	fi

check-vet:
	@echo "==> Checking code issues with go vet:"
	go vet -v ./ 2>&1
	@OUTPUT=$$?; \
	if [ $$OUTPUT -eq 0 ]; then \
		echo "No errors found!"; \
	fi \

check-lint:
	@echo "==> Checking code issues with golint"
	@${GOLINT} 2>&1

fix-format:
	@echo "==> Fixing format with gofmt:"
	@for file in $$(find ./ -iname "*.go"); do \
		echo -n "checking $${file:2}" ; \
		errors=$$(gofmt -d -e -s $$file | grep -c -E "^\+"); \
		if [ "$$errors" -gt 0 ]; then \
			echo " ... fixing $$errors" issues; \
			gofmt -s -w $$file; \
		else \
			echo " ... ok"; \
		fi; \
	done;

test:
	@cd tests; go get
	@ERRORS=0; \
	COUNT=$$(ls tests/*_test.go 2> /dev/null | wc -l); \
	echo -e "\e[32m$$COUNT suites found\e[0m"; \
	for test in $$(ls tests/*_test.go 2> /dev/null) ; do \
		echo -n "==> Running $$test: "; \
		go test $$test; \
		if [ $$? -ne 0 ]; then \
			ERRORS=$$(( ERRORS + 1)); \
		fi; \
	done; \
	if [ $$ERRORS -gt 0 ]; then \
		echo -e "\e[31m$$ERRORS suites failing\e[0m"; \
		exit 1; \
	else \
		echo -e "\e[32mNo errors found!\e[0m"; \
	fi; \

setup:
	go get -u github.com/golang/lint/golint
