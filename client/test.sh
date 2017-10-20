#!/bin/bash
for((i=0; i< 200; i++))
do
	watch -n 1 ./client0 >/dev/null 2>&1 &
	echo "startup client0 ....."$i
#	watch -n 1 ./client1 >/dev/null 2>&1 &
#	echo "startup client1 ....."$i
#        sleep 0.1
done
