mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))

default:
	@echo $(mkfile_dir)
	@echo $(dir $(realpath $(firstword $(mkfile_dir))))


.PHONY: build clean deploy build-generic build-auth
build:
	# dep ensure -v
	make build-generic build-auth

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

build-auth: gomodgen
	cd $(mkfile_dir)functions/auth
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/auth functions/auth/main.go
	cd $(mkfile_dir)


build-generic: gomodgen
	cd $(mkfile_dir)functions/generic_http_events_processor
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/generic_http_events_processor \
		functions/generic_http_events_processor/main.go \
		functions/generic_http_events_processor/signature.go \
		functions/generic_http_events_processor/structs.go 	\
		functions/generic_http_events_processor/utils.go

	cd $(mkfile_dir)
gomodgen:
	chmod u+x $(mkfile_dir)gomod.sh
	$(mkfile_dir)gomod.sh

getPath:
