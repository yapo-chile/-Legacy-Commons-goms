#!/usr/bin/env bash

export MOCK_PORT=8080

file=pact-go_$(uname -s)_amd64.tar.gz

# Validate pact-go binaries
if [ ! -f "$PWD/pact/bin/pact-go" ]; then
  echo "Downloading binaries..."
  echo $PWD
  mkdir -p $PWD/pact/bin
  wget  -P $PWD/pact/bin https://github.com/pact-foundation/pact-go/releases/download/v0.0.11/${file}
  tar zxf $PWD/pact/bin/${file} -C $PWD/pact/bin/
fi

echo "Starting mock response for frontend"

echo -e "\033[0;36m -> open http://localhost:${MOCK_PORT}/api/v1/mock \033[0m"
nohup $PWD/pact/bin/pact/bin/pact-stub-service $PWD/pact/mocks/mock.json --port=${MOCK_PORT} --host=0.0.0.0 --cors
