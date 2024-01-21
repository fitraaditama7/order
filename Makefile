build:
	go build -o ./bin/${APP_NAME} main.go

test:
	go test ./...

test-coverage:
	if [ ! -d "test-coverage" ];then     \
			mkdir test-coverage;           \
	fi
	go test -coverprofile=test-coverage/coverage.out ./... ; go tool cover -func=test-coverage/coverage.out

install-package:
	go mod tidy
	npm --prefix ./view install ./view

gen-mocks:
	mockery --all --output mock --keeptree --disable-version-string

.PHONY: build test test-coverage gen-mocks