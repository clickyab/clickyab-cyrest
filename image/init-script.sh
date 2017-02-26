#!/usr/bin/dumb-init /bin/bash
set -euo pipefail
set -x

MYSQL_USER=${MYSQL_USER:-root}
MYSQL_PASSWORD=${MYSQL_PASSWORD:-bita123}
MYSQL_DB=${MYSQL_DB:-cyrest}
MYSQL_HOST=${MYSQL_HOST:-mysql}
MYSQL_PORT=${MYSQL_PORT:-3306}

AMQP_USER=${AMQP_USER:-cyrest}
AMQP_PASS=${AMQP_PASS:-bita123}
AMQP_HOST=${AMQP_HOST:-rabbitmq}
AMQP_PORT=${AMQP_PORT:-5672}
AMQP_VHOST=${AMQP_VHOST:-tg}

REDIS_PASS=${REDIS_PASS:-}
REDIS_PORT=${REDIS_PORT:-2222}
REDIS_HOST=${REDIS_HOST:-redis}


# TODO : env re-write must be done here
export CYREST_SWAGGER_ROOT=/app/swagger/
export CYREST_FRONT_PATH=/app/public
export CYREST_STATIC_ROOT=/app/statics
export CYREST_SITE=${CYREST_SITE:-rubik.clickyab.ae}
export CYREST_PROTO=${CYREST_PROTO:-http}
export CYREST_REDIS_ADDRESS=${REDIS_HOST}:${REDIS_PORT}
export CYREST_REDIS_PASSWORD=${REDIS_PASS}
export CYREST_MYSQL_DSN="${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/"
export CYREST_MYSQL_DATABASE="${MYSQL_DB}"
export CYREST_PROFILE=disable
export CYREST_SLACK_ACTIVE=true
export CYREST_REDMINE_ACTIVE=true
export CYREST_AMQP_DSN="amqp://${AMQP_USER}:${AMQP_PASS}@${AMQP_HOST}:${AMQP_PORT}/${AMQP_VHOST}"
export CYREST_AMQP_EXCHANGE=${CYREST_AMQP_EXCHANGE:-cy}
export CYREST_TELEGRAM_API_KEY=${CYREST_TELEGRAM_API_KEY:-"347601159:AAEangmt4d67iRwd3-afAaKINzQJKA6q6G4"}
export CYREST_TELEGRAM_BOT_ID=${CYREST_TELEGRAM_BOT_ID:-'$0100000007f9b71480cc4deb7a767921'}
export CYREST_TELEGRAM_BOT_NAME=${CYREST_TELEGRAM_BOT_NAME:-'rubikaddemobot'}
export CYREST_TELEGRAM_CLI_HOST=${CYREST_TELEGRAM_CLI_HOST:-tgcli}
export CYREST_TELEGRAM_CLI_PORT=${CYREST_TELEGRAM_CLI_PORT:-9999}


if [ "$1" = '/app/bin/server' ];
then
	/app/bin/migration -action=up
	exec "$@"
else
	exec "$@"
fi;
