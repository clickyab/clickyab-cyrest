#!/usr/bin/dumb-init /bin/bash
set -euo pipefail
set -x

TG_PORT=${TG_PORT:-9999}

if [ "$1" = 'tg' ];
then
	chown -R telegramd. /home/telegramd
	/tg/bin/telegram-cli -d --json -P ${TG_PORT} --accept-any-tcp
	while :
	do
		echo "Press [CTRL+C] to stop.."
		sleep 10
	done
else
	exec "$@"
fi;