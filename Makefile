.default: test
.phony: test

clean:
	@find build -not -name .gitignore -delete

build: clean
	@go build \
	-ldflags "-X github.com/codereaper/lane/cmd.versionString=0.0.0" \
	-o build/

build-docs: build
	@mkdir build/docs
	@build/lane documentation -o build/docs

tidy: clean
	@go fmt
	@go mod tidy
ifeq ($(strip $(CI)),)
	@git diff --quiet --exit-code || echo 'Warning: Workplace is dirty'
else
	@git diff --quiet --exit-code || (echo 'Error: Workplace is dirty'; exit 1)
endif

unit-tests:
	@go test -timeout 10s ./internal/...

TOOL_VERSION = $(shell grep '^golang ' .tool-versions | sed 's/golang //')
MOD_VERSION = $(shell grep '^go ' go.mod | sed 's/go //')
verify-version:
ifneq ($(TOOL_VERSION),$(MOD_VERSION))
	@echo 'Mismatched go versions'
	@exit 1
endif
	@exit 0

test: verify-version tidy build unit-tests
