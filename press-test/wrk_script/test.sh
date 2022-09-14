#!/bin/sh

basepath=$(cd `dirname $0`; pwd)
cd $basepath

connect_num=$1
stress_time=$2
featureTag=$3
tenant=$4
isUpdate=$5
a=${isUpdate:-true}

wrk --latency -t 8 -c $connect_num -d ${stress_time}s -s script.lua -- "http://10.1.9.188:11149/api-1.0-SNAPSHOT/$tenant/run-rule/compute-features/external/checkpoint?featureTagName=$featureTag&isUpdate=$a"
