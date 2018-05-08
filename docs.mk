export DOCS_DIR=docs
export DOCS_HOST=localhost:$(call genport,3)

## Starts godoc webserver with live docs for the project
docs-start:
	godoc -http ${DOCS_HOST} & echo $$! > .docs.pid

## Stops godoc webserver if running
docs-stop:
	cat .docs.pid | xargs kill

## Compiles static documentation to docs folder
docs-compile: docs-start
	wget -r -np -N -E -p -k -e robots=off --include-directories="/pkg,/lib" --exclude-directories="*" http://${DOCS_HOST}/pkg/github.schibsted.io/Yapo/goms/ || true
	mkdir -p ${DOCS_DIR}
	rm -rf ${DOCS_DIR}/{pkg,lib}
	cp -a ${DOCS_HOST}/{pkg,lib} ${DOCS_DIR}
	rm -rf ${DOCS_HOST}

## Opens the live documentation on the default web browser
docs: docs-start
	open http://${DOCS_HOST}/pkg/github.schibsted.io/Yapo/goms/
