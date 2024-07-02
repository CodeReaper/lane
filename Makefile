.default: test
.phony: test

CI = $(shell env | grep ^CI=)
VERSION = 0.0.0
TOOL_VERSION = $(shell grep '^golang ' .tool-versions | sed 's/golang //')
MOD_VERSION = $(shell grep '^go ' go.mod | sed 's/go //')

clean:
	@find build -not -name .gitignore -delete

build: clean
	go build \
	-trimpath \
	-ldflags "-s -w -X github.com/codereaper/lane/cmd.versionString=$(VERSION)" \
	-o build/bin/

update-docs: build
	@mkdir -p docs/generated
	build/bin/lane documentation -o docs/generated

package: build
	tar -cJvf build/lane-$(VERSION).tar.xz build/bin/
	cd build && sha512sum lane-$(VERSION).tar.xz > lane-$(VERSION).tar.xz.sha512sum

tidy: clean
	go fmt
	go mod tidy
ifeq ($(strip $(CI)),)
	@git diff --quiet --exit-code || echo 'Warning: Workplace is dirty'
else
	@git diff --quiet --exit-code || (echo 'Error: Workplace is dirty'; exit 1)
endif

unit-tests:
	go test -timeout 10s -coverprofile=build/coverage.out ./internal/...
	go tool cover -html=build/coverage.out -o build/coverage.html

verify-version:
ifneq ($(TOOL_VERSION),$(MOD_VERSION))
	@echo 'Mismatched go versions'
	@exit 1
endif
	@exit 0

test: verify-version tidy unit-tests
