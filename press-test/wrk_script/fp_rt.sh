#!/bin/sh

a=10
while [ "$a" -lt 120 ]
do
   wrk --latency -t 1 -c $a -d 1m -s script.lua -- "http://k8s-us-dev.dv-api.com" "/yuga1/fp-rt/api-1.0-SNAPSHOT/yb6/run-rule/compute-features/external/rule-set?setId=2&featureTagName=EXPENSIVE_COMPUTE"
   a=`expr $a + 10`
done
