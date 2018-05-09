## Starts godoc webserver with live docs for the project
docs-start:
	godoc -http ${DOCS_HOST} &> /dev/null & echo $$! > .docs.pid

## Stops godoc webserver if running
docs-stop:
	cat .docs.pid | xargs kill

## Compiles static documentation to docs folder
docs-compile: docs-start
	@scripts/commands/docs-compile.sh

## Generates a commit updating the docs
docs-update: docs-compile
	git add docs
	git commit -m "${DOCS_COMMIT_MESSAGE}"

## Opens the live documentation on the default web browser
docs: docs-start
	open http://${DOCS_HOST}/pkg/github.schibsted.io/Yapo/goms/
