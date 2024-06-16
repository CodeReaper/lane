.default: test

clean:
	@find build -not -name .gitignore -delete

build: clean
	@go build \
	-ldflags "-X github.com/codereaper/lane/cmd.versionString=0.0.0" \
	-o build/

build-docs: build
	@mkdir build/docs
	@build/lane documentation -o build/docs

test-clean: clean
	@go fmt
	@go mod tidy
ifeq ($(strip $(CI)),)
	@git diff --quiet --exit-code || echo 'Workplace is dirty'
else
	@git diff --quiet --exit-code || (echo 'Workplace is dirty'; exit 1)
endif

test: test-clean build
	@go test -timeout 10s ./internal/...
