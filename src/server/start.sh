#!/bin/bash

./server -name=GateServer -index=0 &

./server -name=LoginServer -index=0 &
./server -name=LoginServer -index=1 &
