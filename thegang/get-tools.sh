#!/bin/bash

xcode-select --install
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
brew update
brew install mercurial
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
gvm install go1.20
gvm use go1.20
brew install golangci-lint
brew install pre-commit
pre-commit install
cp ./thegang/docker-compose.yaml .
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
go install -v github.com/go-critic/go-critic/cmd/gocritic@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golang/mock/mockgen@v1.6.0
go get github.com/google/wire/cmd/wire@v0.5.0
go install github.com/google/wire/cmd/wire@latest
go install github.com/google/wire/cmd/wire
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
go install github.com/pressly/goose/v3/cmd/goose@latest