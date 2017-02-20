.PHONY : clean build rpm-clean rpm-build
get_pid=$$(ps ux | grep -v grep | grep -E "\./goms" | awk '{ print $$2 }')

EXEC="./goms"
GOLINT=$$GOPATH/bin/golint

GENPORTOFF?=0
genport = $(shell expr ${GENPORTOFF} + \( $(shell id -u) - \( $(shell id -u) / 100 \) \* 100 \) \* 200 + 30100 + $(1))

ifndef ENVIRONMENT
	ENVIRONMENT=Develop
	YO=`whoami`
	LISTEN_PORT=$(call genport,2)
	SERVER_ROOT=${PWD}
	DOCUMENT_ROOT=${PWD}/src/public
	SERVERNAME=`hostname`
	BASE_URL="http://"${SERVERNAME}":"${LISTEN_PORT}
endif

info:
	@echo "YO           : "${YO}
	@echo "ServerRoot   : "${SERVER_ROOT}
	@echo "DocumentRoot : "${DOCUMENT_ROOT}
	@echo "API Base URL : ${BASE_URL}"

build:
	@${MAKE} stop
	go get
	go build

run:
	PID=$(call get_pid); \
	if [ "$$PID" ]; then \
		echo "ERROR: service is already running (PID: $$PID)"; \
	else if [ ! -x "${EXEC}" ]; then \
		echo "service ${EXEC} not found"; \
	else \
		${MAKE} update_config; \
		${MAKE} demonize; \
	fi; \
	fi

start:
	PID=$(call get_pid); \
	if [ "$$PID" ]; then \
		echo "ERROR: service is already running (PID: $$PID)"; \
	else \
		${MAKE} build; \
		${MAKE} run; \
	fi

stop:
	@PID=$(call get_pid); \
	if [ "$$PID" ]; then \
		kill -15 $$PID; \
	fi

restart:
	${MAKE} stop
	${MAKE} start

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
	@nohup ${EXEC} >> ${EXEC}.log 2>&1 &

rpm-clean:
	PID=$(call get_pid); \
	if [ "$$PID" ]; then \
		echo "ERROR: service is running (PID: $$PID)"; \
		exit -1; \
	else \
		rm -Rf build; \
	fi

rpm-build: rpm-clean rpm-setuptree build
	cp conf/conf.json.in conf/conf.json
	rpmbuild -bb goms.spec
	mv build/RPMS/x86_64/goms*.rpm ./
	rm -Rf build

rpm-setuptree:
	mkdir -p build/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

check:
	${MAKE} check-format
	${MAKE} check-vet
	${MAKE} check-lint

check-format:
	echo "==> Checking format with gofmt:"
	ERRORS=0; \
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
	echo "==> Checking code issues with go vet:"
	go vet -v ./ 2>&1
	OUTPUT=$$?; \
	if [ $$OUTPUT -eq 0 ]; then \
		echo "No errors found!"; \
	fi \

check-lint:
	echo "==> Checking code issues with golint"
	${GOLINT} 2>&1

fix-format:
	echo "==> Fixing format with gofmt:"
	for file in $$(find ./ -iname "*.go"); do \
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
	#@${MAKE} db-load-test
	ERRORS=0; \
	for test in tests/*_test.go ; do \
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
		echo "No errors found!"; \
	fi; \

setup:
	go get -u github.com/golang/lint/golint
