#!/bin/bash

baseDir=`pwd`


for ((i = 32; i <= 100; i = i + 2)); do
    wrk --latency -s $baseDir/lua/script-zxb.lua "http://176.34.196.137:31403" -t $i -c 864 -d 30s -- $baseDir/data/du.json test true 
done