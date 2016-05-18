#@IgnoreInspection BashAddShebang
export DEFAULT_PASS=bita123
export GO=$(shell which go)
export GIT:=$(shell which git)
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export BIN=$(ROOT)/bin
export GOPATH=$(ROOT):$(ROOT)/vendor
export WATCH?=hello
export LONGHASH=$(shell git log -n1 --pretty="format:%H" | cat)
export SHORTHASH=$(shell git log -n1 --pretty="format:%h"| cat)
export COMMITDATE=$(shell git log -n1 --pretty="format:%cd"| sed -e "s/ /-/g")
export COMMITCOUNT=$(shell git rev-list HEAD --count| cat)
export BUILDDATE=$(shell date| sed -e "s/ /-/g")
export FLAGS="-X shared/config.hash=$(LONGHASH) -X shared/config.short=$(SHORTHASH) -X shared/config.date=$(COMMITDATE) -X shared/config.count=$(COMMITCOUNT) -X shared/config.build=$(BUILDDATE)"
export LDARG=-ldflags $(FLAGS)
export DBPASS?=$(DEFAULT_PASS)
export POSTGRES_USER?=helium
export RUSER?=helium
export RPASS?=$(DEFAULT_PASS)


.PHONY: all gb clean tools-fswatch hello

all:  $(BIN)/gb
	$(BIN)/gb build $(LDARG)

needroot :
	@[ "$(shell id -u)" -eq "0" ] || exit 1

notroot :
	@[ "$(shell id -u)" != "0" ] || exit 1

gb: notroot
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(GO) get -u -v github.com/constabulary/gb/...

clean:
	rm -rf $(ROOT)/pkg $(ROOT)/vendor/pkg
	cd $(ROOT) && git clean -fX ./bin

$(BIN)/gb: notroot
	@[ -f $(BIN)/gb ] || make gb

#
# Tools
#
tools-fswatch: $(BIN)/gb
	$(BIN)/gb build $(LDARG) tools/fswatch

tools-godebug: $(BIN)/gb
	$(BIN)/gb build $(LDARG) github.com/mailgun/godebug

tools-golint: $(BIN)/gb
	$(BIN)/gb build $(LDARG) github.com/golang/lint/golint

tools-govet: $(BIN)/gb
	$(BIN)/gb build $(LDARG) golang.org/x/tools/cmd/vet

tools-migrate: $(BIN)/gb
	$(BIN)/gb build $(LDARG) apps/migrate

tools-gerrithook: $(BIN)/gb
	$(BIN)/gb build $(LDARG) tools/gerrithook

tools-goimports: $(BIN)/gb
	$(BIN)/gb build $(LDARG) golang.org/x/tools/cmd/goimports

tools-gotype: $(BIN)/gb
	$(BIN)/gb build $(LDARG) golang.org/x/tools/cmd/gotype

tools-godoc: $(BIN)/gb
	$(BIN)/gb build $(LDARG) golang.org/x/tools/cmd/godoc

tools-fgt: $(BIN)/gb
	$(BIN)/gb build $(LDARG) tools/fgt

tools-deadcode: $(BIN)/gb
	$(BIN)/gb build $(LDARG) tools/deadcode

tools-ineffassign: $(BIN)/gb
	$(BIN)/gb build $(LDARG) tools/ineffassign

tools-goconvey: $(BIN)/gb
	$(BIN)/gb build $(LDARG) github.com/smartystreets/goconvey/

tools-codegen: $(BIN)/gb
	$(BIN)/gb build $(LDARG) tools/codegen

tools-errcheck: $(BIN)/gb
	$(BIN)/gb build $(LDARG) github.com/kisielk/errcheck

godoc: tools-godoc
	#open localhost:6060 for doc, Ctrl+C to stop
	$(BIN)/godoc -http=:6060

errcheck: tools-errcheck
	find ./src/apps/* -type d | sed 's|./src/||' | grep -v sdpctl | xargs $(BIN)/errcheck

protobuf-go:
	$(BIN)/gb build $(LDARG) github.com/golang/protobuf/protoc-gen-go


protobuf: notroot
	wget -c -O $(ROOT)/tmp/protobuf-beta3.zip https://github.com/google/protobuf/archive/v3.0.0-beta-3.zip
	cd $(ROOT)/tmp && unzip -o $(ROOT)/tmp/protobuf-beta3.zip
	cd $(ROOT)/tmp/protobuf-3.0.0-beta-3/ && ./autogen.sh
	cd $(ROOT)/tmp/protobuf-3.0.0-beta-3/ && ./configure --prefix=/usr
	cd $(ROOT)/tmp/protobuf-3.0.0-beta-3/ && make

install-protobuf: needroot
	cd $(ROOT)/tmp/protobuf-3.0.0-beta-3/ && make install
#
# Migration
#

migup: tools-migrate
	$(BIN)/migrate --action=up --folder=$(ROOT)/db/migrations

migdown: tools-migrate
	$(BIN)/migrate --action=down --folder=$(ROOT)/db/migrations

migdown-all: tools-migrate
	$(BIN)/migrate --action=down-all --folder=$(ROOT)/db/migrations

migredo: tools-migrate
	$(BIN)/migrate --action=redo --folder=$(ROOT)/db/migrations

miglist: tools-migrate
	$(BIN)/migrate --action=list --folder=$(ROOT)/db/migrations

migcreate:
	/bin/bash $(BIN)/create_migration.sh

goimports: tools-goimports
	$(BIN)/goimports -w $(ROOT)/src

watch: $(WATCH) tools-fswatch
	$(BIN)/fswatch -d 10 -ext go make run-$(WATCH)

build-protobuf:
	protoc --go_out $(ROOT)/src/gen -I $(ROOT)/src/proto/ $(ROOT)/src/proto/*.proto

#	protoc --java_out $(ROOT)/src/dummy -I $(ROOT)/src/proto/ $(ROOT)/src/proto/*.proto

#
# Codegen
#

codegen-user: tools-codegen
	@$(BIN)/codegen -p modules/user/controllers
	@$(BIN)/codegen -p modules/user/aaa


codegen-audit: tools-codegen
	@$(BIN)/codegen -p modules/audit/controllers

codegen-misc: tools-codegen
	@$(BIN)/codegen -p modules/misc/controllers
	@$(BIN)/codegen -p modules/misc/t9n

codegen-balance: tools-codegen
	@$(BIN)/codegen -p modules/balance/acc
	@$(BIN)/codegen -p modules/balance/controllers

swagger-cleaner:
	@rm -f $(WORK_DIR)/swagger/*.json
	@rm -f $(WORK_DIR)/swagger/*.yaml

swagger-client: tools-swagger
	GOPATH=$(ROOT) cd $(ROOT)/src && $(BIN)/swagger generate client -f $(ROOT)/3rd/swagger/phoenix.yaml

codegen: swagger-cleaner codegen-user codegen-audit codegen-balance codegen-misc
	@cp $(WORK_DIR)/swagger/phoenix.yaml $(ROOT)/3rd/swagger
	@cp $(WORK_DIR)/swagger/phoenix.json $(ROOT)/3rd/swagger
	@echo "Done"

#
# Lint
#

vet: tools-govet tools-fgt
	@$(BIN)/fgt $(BIN)/vet $(ROOT)/src
	@$(BIN)/fgt $(BIN)/vet --shadow $(ROOT)/src
	@echo "vet is done"

golint: tools-golint tools-fgt
	@$(BIN)/fgt $(BIN)/golint $(ROOT)/src
	@echo "lint is done"

gotype: tools-gotype tools-fgt
	@$(BIN)/fgt $(BIN)/gotype $(ROOT)/src
	@echo "type is done"

ineffassign: tools-ineffassign tools-fgt
	@$(BIN)/fgt $(BIN)/ineffassign $(ROOT)/src
	@echo "inefassign is done"

lint: goimports all errcheck vet golint gotype ineffassign
	@echo "All done"

postgres-setup: needroot
	sudo -u postgres psql -U postgres -d postgres -c "CREATE USER $(POSTGRES_USER) WITH PASSWORD '$(DBPASS)';" || sudo -u postgres psql -U postgres -d postgres -c "ALTER USER $(POSTGRES_USER) WITH PASSWORD '$(DBPASS)';"
	sudo -u postgres psql -U postgres -c "CREATE DATABASE helium;" || echo "Database helium is already there?"
	sudo -u postgres psql -U postgres -c "GRANT ALL ON DATABASE helium TO $(POSTGRES_USER);"