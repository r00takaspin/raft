#!/bin/bash

killall raft
make

./raft -p 10001 -nodes Ironclad.local:10002,Ironclad.local:10003 &
./raft -p 10002 -nodes Ironclad.local:10001,Ironclad.local:10003 &
./raft -p 10003 -nodes Ironclad.local:10001,Ironclad.local:10002 &