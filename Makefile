#@IgnoreInspection BashAddShebang
export APPNAME=cyrest
export DEFAULT_PASS=bita123
export GO=$(shell which go)
export GIT:=$(shell which git)
export DIFF:=$(shell which diff)
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export BIN=$(ROOT)/bin
export GB=$(BIN)/gb
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
export DB_HOST?=127.0.0.1
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
	$(BUILD)

needroot :
	@[ "$(shell id -u)" -eq "0" ] || exit 1

gb:
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(GO) get $(UPDATE) -v github.com/constabulary/gb/...

clean:
	rm -rf $(ROOT)/pkg $(ROOT)/vendor/pkg
	cd $(ROOT) && git clean -fX ./bin
	@echo "Done"

$(BIN)/gb:
	[ -f $(BIN)/gb ] || make gb


server:
	$(BUILD) server

run-server: server
	sudo setcap cap_net_bind_service=+ep $(BIN)/server
	$(BIN)/server

watch-server: codegen
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
	cd $(ROOT) && $(BIN)/go-bindata -o ./src/migration/migration.gen.go -nomemcopy=true -pkg=main ./db/migrations/...

tools-migrate: $(BIN)/gb migration
	$(BUILD) migration

watch: $(WATCH) tools-fswatch
	$(BIN)/fswatch -d 10 -ext go make run-$(WATCH)

#
# Swagger
#

swagger-ui:
	$(GIT) clone --depth 1 https://github.com/swagger-api/swagger-ui.git $(ROOT)/tmp/swagger-ui || true
	cp -R $(ROOT)/tmp/swagger-ui/dist/* $(ROOT)/3rd/swagger


#
# Codegen
#

codegen-user: tools-codegen
	$(BIN)/codegen -p modules/user/controllers
	$(BIN)/codegen -p modules/user/aaa

codegen-category: tools-codegen
	@$(BIN)/codegen -p modules/category/controllers
	@$(BIN)/codegen -p modules/category/cat

codegen-misc: tools-codegen
	@$(BIN)/codegen -p modules/misc/base
	@$(BIN)/codegen -p modules/misc/controllers
	@$(BIN)/codegen -p modules/misc/t9n

codegen-channel: tools-codegen
	$(BIN)/codegen -p modules/channel/controllers
	$(BIN)/codegen -p modules/channel/chn

swagger-cleaner:
	@rm -f $(WORK_DIR)/swagger/*.json
	@rm -f $(WORK_DIR)/swagger/*.yaml

swagger-client: tools-swagger
	GOPATH=$(ROOT) cd $(ROOT)/src && $(BIN)/swagger generate client -f $(ROOT)/3rd/swagger/cyrest.yaml

codegen: swagger-ui swagger-cleaner codegen-misc codegen-user codegen-category codegen-channel
	@cp $(WORK_DIR)/swagger/out.yaml $(ROOT)/3rd/swagger/cyrest.yaml
	@cp $(WORK_DIR)/swagger/out.json $(ROOT)/3rd/swagger/cyrest.json
	@echo "Done"

#
# Lint
#
metalinter:
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(GO) get $(UPDATE) -v github.com/alecthomas/gometalinter
	GOPATH=$(ROOT)/tmp GOBIN=$(ROOT)/bin $(ROOT)/bin/gometalinter --install

$(BIN)/gometalinter:
	@[ -f $(BIN)/gometalinter ] || make metalinter

lint-common: $(BIN)/gometalinter
	$(LINTER) $(ROOT)/src/common/...

lint-modules: $(BIN)/gometalinter
	$(LINTER) $(ROOT)/src/modules/...

lint-mains: $(BIN)/gometalinter
	$(LINTER) $(ROOT)/src/server/...

lint: lint-common lint-modules lint-mains
	@echo "Done"

mysql-createdb:
	echo 'DROP DATABASE IF EXISTS cyrest;' | mysql -h $(DB_HOST) -u $(DB_USER) -p$(DBPASS)
	echo 'CREATE DATABASE cyrest;' | mysql -h $(DB_HOST) -u $(DB_USER) -p$(DBPASS)

mysql-setup: needroot
	echo 'UPDATE user SET plugin="";' | mysql mysql
	echo 'UPDATE user SET password=PASSWORD("$(DBPASS)") WHERE user="$(DB_USER)";' | mysql mysql
	echo 'FLUSH PRIVILEGES;' | mysql mysql
	make mysql-createdb


setcap: $(BIN)/server needroot
	setcap cap_net_bind_service=+ep $(BIN)/server

restore: $(GB)
	PATH=$(PATH):$(BIN) $(GB) vendor restore
	cp $(ROOT)/vendor/manifest $(ROOT)/vendor/manifest.done

conditional-restore:
	$(DIFF) $(ROOT)/vendor/manifest $(ROOT)/vendor/manifest.done || make restore

docker-build: conditional-restore codegen all
