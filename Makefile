basedir = $(shell pwd)
gopath = "$(basedir)/vendor:$(GOPATH)"

.PNONY: all test deps fmt clean check-gopath

all: check-gopath clean fmt deps test
	@echo "==> Compiling source code."
	@env GOPATH=$(gopath) go build -v -o ./bin/cfn-clone ./cfn-clone

race: check-gopath clean fmt deps test
	@echo "==> Compiling source code with race detection enabled."
	@env GOPATH=$(gopath) go build -race -o ./bin/cfn-clone ./cfn-clone

test: check-gopath
	@echo "==> Running tests."
	@env GOPATH=$(gopath) go test -cover ./cfn-clone

deps: check-gopath
	@echo "==> Downloading dependencies."
	@env GOPATH=$(gopath) go get -d -v ./cfn-clone/...
	@echo "==> Removing SCM files from vendor."
	@find ./vendor -type d -name .git | xargs rm -rf
	@find ./vendor -type d -name .bzr | xargs rm -rf
	@find ./vendor -type d -name .hg | xargs rm -rf

fmt:
	@echo "==> Formatting source code."
	@gofmt -w ./cfn-clone

clean:
	@echo "==> Cleaning up previous builds."
	@rm -rf "$(GOPATH)/pkg" ./vendor/pkg ./bin

check-gopath:
ifndef GOPATH
	$(error GOPATH is undefined)
endif
