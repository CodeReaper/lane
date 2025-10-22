# cSpell:ignore trimpath coverprofile

CI = $(shell env | grep ^CI=)
VERSION = 0.0.0
SUFFIX =
TOOL_VERSION = $(shell grep '^golang ' .tool-versions | sed 's/golang //')
MOD_VERSION = $(shell grep '^go ' go.mod | sed 's/go //')

clean:
	@find build -type f ! -name .gitignore -exec rm {} +
	@find build -type d -mindepth 1 -exec rmdir {} +

build: clean
	go build \
	-trimpath \
	-ldflags "-s -w -X github.com/codereaper/lane/cmd.versionString=$(VERSION)" \
	-o build/bin/
	cp LICENSE build/bin/LICENSE.txt

update-docs: build
	@mkdir -p docs/generated
	build/bin/lane documentation -o docs/generated

package: build
	cd build/bin && tar -cJvf ../lane-$(VERSION)$(SUFFIX).tar.xz *
	cd build && sha512sum lane-$(VERSION)$(SUFFIX).tar.xz > lane-$(VERSION)$(SUFFIX).tar.xz.sha512sum

tidy: clean
	go fmt
	go mod tidy
ifeq ($(strip $(CI)),)
	@git diff --quiet --exit-code || echo 'Warning: Workplace is dirty'
else
	@git diff --quiet --exit-code || (echo 'Error: Workplace is dirty'; exit 1)
endif

unit-tests:
	go test -timeout 10s -p 1 -coverprofile=build/coverage.out ./internal/...
ifeq ($(strip $(GITHUB_STEP_SUMMARY)),)
	go tool cover -html=build/coverage.out -o build/coverage.html
	go tool cover -func=build/coverage.out
else
	{	\
		echo '## Code Coverage'; \
		echo '|File|Method|Coverage|'; \
		echo '|---|---|--:|'; \
		go tool cover -func=build/coverage.out | awk -F'\t+' 'BEGIN {OFS=" | "} {print "| " $$1, $$2, $$3 " |"}'; \
	} | tee -a $(GITHUB_STEP_SUMMARY)
endif

smoke-test:
	@if [ -f credentials.json ]; then \
		go run ./... translations list -c credentials.json; \
	fi

verify-version:
ifneq ($(TOOL_VERSION),$(MOD_VERSION))
	@echo 'Mismatched go versions'
	@exit 1
endif

test: verify-version tidy unit-tests smoke-test
