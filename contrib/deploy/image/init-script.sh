#!/usr/bin/dumb-init /bin/bash
set -euo pipefail
set -x

MYSQL_USER=${MY_USER:-root}
MYSQL_PASSWORD=${MY_PASS:-bita123}
MYSQL_DB=${MY_DB:-cyrest}

# TODO : env re-write must be done here
export CYREST_SWAGGER_ROOT=/app/swagger/
export CYREST_SITE=rubik.clickyab.ae
export CYREST_PROTO=http
export CYREST_REDIS_ADDRESS=redis:6379
export CYREST_MYSQL_DSN="${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(mysql:3306)/"
export CYREST_MYSQL_DATABASE="${MYSQL_DB}"
export CYREST_PROFILE=disable

if [ "$1" = '/app/bin/server' ];
then
	/app/bin/migration -action=up
	exec "$@"
else
	exec "$@"
fi;