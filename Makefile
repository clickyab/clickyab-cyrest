#@IgnoreInspection BashAddShebang
export APPNAME=cyrest
export DEFAULT_PASS=bita123
export GO=$(shell which go)
export GIT:=$(shell which git)
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export BIN=$(ROOT)/bin
export GOPATH=$(ROOT):$(ROOT)/vendor
export WATCH?=hello
export LONGHASH=$(shell git log -n1 --pretty="format:%H" | cat)
export SHORTHASH=$(shell git log -n1 --pretty="format:%h"| cat)
export COMMITDATE=$(shell git log -n1 --date="format:%D-%H-%I-%S" --pretty="format:%cd"| sed -e "s/\//-/g")
export COMMITCOUNT=$(shell git rev-list HEAD --count| cat)
export BUILDDATE=$(shell date "+%D/%H/%I/%S"| sed -e "s/\//-/g")
export FLAGS="-X common/version.hash=$(LONGHASH) -X common/version.short=$(SHORTHASH) -X common/version.date=$(COMMITDATE) -X common/version.count=$(COMMITCOUNT) -X common/version.build=$(BUILDDATE)"
export LDARG=-ldflags $(FLAGS)
export BUILD=$(BIN)/gb build $(LDARG)
export DBPASS?=$(DEFAULT_PASS)
export DB_USER?=root
export RUSER?=$(APPNAME)
export RPASS?=$(DEFAULT_PASS)
export WORK_DIR=$(ROOT)/tmp
export LINTER=$(BIN)/gometalinter -e ".*src/modules/user/templates/mail.go.*" --cyclo-over=15 --line-length=120 --deadline=100s --disable-all --enable=structcheck --enable=aligncheck --enable=deadcode --enable=gocyclo --enable=ineffassign --enable=golint --enable=goimports --enable=errcheck --enable=varcheck --enable=interfacer --enable=goconst --enable=gosimple --enable=staticcheck --enable=unused --enable=misspell --enable=lll

ifdef UPDATEGB
export UPDATE=-u
else
export UPDATE=
endif


all:  $(BIN)/gb
	$(BIN)/gb build $(LDARG)

needroot :
	@[ "$(shell id -u)" -eq "0" ] || exit 1

notroot :
	@[ "$(shell id -u)" != "0" ] || exit 1

gb: notroot
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(GO) get $(UPDATE) -v github.com/constabulary/gb/...

clean:
	rm -rf $(ROOT)/pkg $(ROOT)/vendor/pkg
	cd $(ROOT) && git clean -fX ./bin
	@echo "Done"

$(BIN)/gb: notroot
	[ -f $(BIN)/gb ] || make gb


server:
	$(BUILD) server

run-server: server
	sudo setcap cap_net_bind_service=+ep $(BIN)/server
	$(BIN)/server

watch-server:
	make watch WATCH=server

#
# Tools
#
tools-fswatch: $(BIN)/gb
	$(BIN)/gb build $(LDARG) tools/fswatch

tools-codegen: $(BIN)/gb
	$(BIN)/gb build $(LDARG) tools/codegen

tools-gobindata: $(BIN)/gb
	$(BIN)/gb build $(LDARG) github.com/jteeuwen/go-bindata/go-bindata

godoc: tools-godoc
	#open localhost:6060 for doc, Ctrl+C to stop
	$(BIN)/godoc -http=:6060

errcheck: tools-errcheck
	find ./src/apps/* -type d | sed 's|./src/||' | grep -v sdpctl | xargs $(BIN)/errcheck

#
# Migration
#

migup: tools-migrate
	$(BIN)/migration -action=up

migdown: tools-migrate
	$(BIN)/migration -action=down

migdown-all: tools-migrate
	$(BIN)/migration -action=down-all

migredo: tools-migrate
	$(BIN)/migration -action=redo

miglist: tools-migrate
	$(BIN)/migration -action=list

migcreate:
	@/bin/bash $(BIN)/create_migration.sh

migcp:
	mkdir -p $(ROOT)/db/migrations
	cd $(ROOT)/db/migrations && rm -f *.sql
	cp $(ROOT)/src/modules/*/migrations/*.sql $(ROOT)/db/migrations

migration: migcp tools-gobindata
	cd $(ROOT) && $(BIN)/go-bindata -o ./src/migration/migration.go -nomemcopy=true -pkg=main ./db/migrations/...

tools-migrate: $(BIN)/gb migration
	$(BUILD) migration

watch: $(WATCH) tools-fswatch
	$(BIN)/fswatch -d 10 -ext go make run-$(WATCH)

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

swagger-cleaner:
	@rm -f $(WORK_DIR)/swagger/*.json
	@rm -f $(WORK_DIR)/swagger/*.yaml

swagger-client: tools-swagger
	GOPATH=$(ROOT) cd $(ROOT)/src && $(BIN)/swagger generate client -f $(ROOT)/3rd/swagger/cyrest.yaml

codegen: swagger-cleaner codegen-user codegen-audit codegen-misc
	@cp $(WORK_DIR)/swagger/cyrest.yaml $(ROOT)/3rd/swagger
	@cp $(WORK_DIR)/swagger/cyrest.json $(ROOT)/3rd/swagger
	@echo "Done"

#
# Lint
#
metalinter: notroot
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(GO) get $(UPDATE) -v github.com/alecthomas/gometalinter
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(ROOT)/bin/gometalinter --install

$(BIN)/gometalinter: notroot
	@[ -f $(BIN)/gometalinter ] || make metalinter

lint-common: $(BIN)/gometalinter
	$(LINTER) $(ROOT)/src/common/...

lint-modules: $(BIN)/gometalinter
	$(LINTER) $(ROOT)/src/modules/...

lint-mains: $(BIN)/gometalinter
	$(LINTER) $(ROOT)/src/server/...

lint: lint-common lint-modules lint-mains
	@echo "Done"


mysql-setup: needroot
	echo 'UPDATE user SET plugin="";' | mysql mysql
	echo 'UPDATE user SET password=PASSWORD("$(DBPASS)") WHERE user="$(DB_USER)";' | mysql mysql
	echo 'FLUSH PRIVILEGES;' | mysql mysql
	echo 'CREATE DATABASE cyrest;' | mysql -u $(DB_USER) -p$(DBPASS)

setcap: $(BIN)/server needroot
	setcap cap_net_bind_service=+ep $(BIN)/server