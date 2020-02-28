Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X github.com/eosio-enterprise/chappe/cmd.Version=$(Version) -X github.com/eosio-enterprise/chappe/cmd.GitCommit=$(GitCommit)"
PLATFORM := $(shell ./hack/platform-tag.sh)
SOURCE_DIRS = cmd pkg main.go

.PHONY: all

.PHONY: build
build:
	go build

.PHONY: fmtcheck
fmt: ## Checks for style violation using gofmt
	@gofmt -l -s $(SOURCE_DIRS) | grep ".*\.go"; if [ "$$?" = "0" ]; then echo "Run gofmt -w -s ./path" && exit 1; fi

.PHONY: test
test:
	CGO_ENABLED=0 go test $(shell go list ./... | grep -v /vendor/|xargs echo) -cover

.PHONY: dist
dist:
	mkdir -p bin
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/chappe
	GO111MODULE=on CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/chappe-darwin
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/chappe-armhf
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/chappe-arm64
	GO111MODULE=on CGO_ENABLED=0 GOOS=windows go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/chappe.exe
