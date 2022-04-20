TEST?=$$(go list ./citrixadm/... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=citrix
NAME=citrixadm
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=darwin_amd64

default: install

docgen:
	tfplugindocs

build: fmt
	go build -o ${BINARY}

debug-build: fmt
	go build -gcflags="all=-N -l" -o ${BINARY}
	cp -f ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

fmt:
	go fmt ./...

tffmt:
	terraform fmt -list=true -recursive examples


install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

testacc:
	# Usage: make testacc VPX_IP=10.0.1.76 VPX_USER=nsroot VPX_PASSWORD=verysecretpassword AGENT_IP=10.0.1.91
	rm -i citrixadm/citrixadm.acctest.log
	TF_ACC=1 TF_ACC_LOG_PATH=./citrixadm.acctest.log TF_LOG=TRACE VPX_IP=$(VPX_IP) VPX_USER=$(VPX_USER) VPX_PASSWORD=$(VPX_PASSWORD) AGENT_IP=$(AGENT_IP) go test terraform-provider-citrixadm/citrixadm -v

start-debug: debug-build
	~/go/bin/dlv exec --accept-multiclient --continue --headless ./${BINARY} -- -debug
