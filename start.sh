#!/bin/bash

trap "cleanup" SIGINT SIGTERM
echo "pid is $$"


cleanup(){
  echo "---> Caught signals and Let us do clean-ups"
  kill -9 $OFFLINE_PID
  kill -9 $SERVER_PID
  exit 2
}

# python3 offline.py >/dev/null 2>&1 &
python3 offline.py &
OFFLINE_PID=$!

while true
do
	python3 server.py &
	SERVER_PID=$!

	sleep 30m
	kill -9 $SERVER_PID
done