#!/bin/bash
/etc/init.d/ssh start
tmass load --layout-dir=/home/develop/helium/bin/ init.yml

# DO NOT ALTER THE FOLLOWING LINE
while true;
do
	echo "Docker main shell, press Ctrl+C to kill this container"
	sleep 1
done;
