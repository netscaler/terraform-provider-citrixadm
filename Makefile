TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=citrix
NAME=citrixadm
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=darwin_amd64

default: install

build: fmt
	go build -o ${BINARY}

debug-build: fmt
	go build -gcflags="all=-N -l" -o ${BINARY}
	cp -f ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

fmt:
	go fmt ./...

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

start-debug: debug-build
	~/go/bin/dlv exec --accept-multiclient --continue --headless ./${BINARY} -- -debug
