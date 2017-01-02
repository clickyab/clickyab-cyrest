#!/usr/bin/env bash

function usage(){
	echo "${0} USERNAME PORT"
	exit 1
}

function detect(){
  type -P $1  || { echo "Require $1 but not installed. Aborting." >&2; exit 1; }
}

USER=${1:-}
PORT=${2:-9000}
CLI_PATH=${3:-telegram-cli}

detect ${CLI_PATH}

if [ -z ${USER} ]; then
	usage
fi;

if [ ! $(id -u ${USER}) ]; then
	if [ $(useradd -mr ${USER}) ];then
		echo "can not create the user ${USER}"
		exit 2
	fi;
fi

su ${USER} -c "${CLI_PATH} -P ${PORT} -d --json"
