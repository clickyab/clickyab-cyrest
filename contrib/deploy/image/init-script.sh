#!/usr/bin/dumb-init /bin/bash
set -euo pipefail
set -x

MYSQL_USER=${MY_USER:-root}
MYSQL_PASSWORD=${MY_PASS:-bita123}
MYSQL_DB=${MY_DB:-cyrest}

# TODO : env re-write must be done here
export CYREST_SWAGGER_ROOT=/app/swagger/
export CYREST_FRONT_PATH=/app/public
export CYREST_SITE=rubik.clickyab.ae
export CYREST_PROTO=http
export CYREST_REDIS_ADDRESS=redis:6379
export CYREST_MYSQL_DSN="${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(mysql:3306)/"
export CYREST_MYSQL_DATABASE="${MYSQL_DB}"
export CYREST_PROFILE=disable
export CYREST_SLACK_ACTIVE=true
export CYREST_REDMINE_ACTIVE=true
export CYREST_AMQP_DSN="amqp://cyrest:bita123@rabbitmq:5672/"
export CYREST_AMQP_EXCHANGE="cy"
export CYREST_TELEGRAM_API_KEY="273335144:AAEv4uPeo68X7Scc3MLKxwMO1YI3JFkWiJM"
export CYREST_TELEGRAM_BOT_ID='$0100000068c34a10ed72226be64e8d4d'
export CYREST_TELEGRAM_CLI_HOST=tgcli
export CYREST_TELEGRAM_CLI_PORT=9999

if [ "$1" = '/app/bin/server' ];
then
	/app/bin/migration -action=up
	exec "$@"
else
	exec "$@"
fi;