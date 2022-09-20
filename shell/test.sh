#!/bin/bash

export KUBECONFIG=~/.ssh/188-kubeconfig

baseDir=`pwd`
replica_fp() {
    kubectl -n workshop scale deploy fp-deployment --replicas $1
    while :; do
        echo "curr time" $(date '+%Y-%m-%d %H:%M:%S')
        kubectl -n gs-test get po

        fp_count=$(kubectl -n workshop get po | grep -c -E "fp.* 1/1")
        if [ $fp_count -eq $1 ]; then
            return
        else
            sleep 30s
        fi

    done
}

alarm() {
    curl https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=9e380bb9-e88d-4e29-a79c-ef6ce98c1aac -H "'Content-Type: application/json'" -d '{"msgtype":"text","text":{"content":"'$1'"}}'
}


wait_kafka_consumer_0() {
    while :; do
        consumer_num=$(kubectl exec -it -n sandbox kafka-0 -- kafka-consumer-groups --bootstrap-server localhost:9092 --group velocity_gs_test --describe | grep "gs_test_fp_velocity" | awk '{if ($5 != 0) {print $5}}' | wc -l)
        kubectl exec -it -n sandbox kafka-0 -- kafka-consumer-groups --bootstrap-server localhost:9092 --group velocity_gs_test --describe | grep "gs_test_fp_velocity"
        echo "consumer_num" $consumer_num
        if [ $consumer_num -eq 0 ]; then
            return
        fi
        sleep 1s
    done
}

replica_cassandra() {
    kubectl -n workshop scale statefulset cassandra --replicas $1
    while :; do
        echo "curr time" $(date '+%Y-%m-%d %H:%M:%S')
        kubectl -n workshop get po

        cassandra_count=$(kubectl -n workshop get po | grep -c -E "cassandra.* 1/1")
        #cassandra_count=$(kubectl -n gs-test get po | grep -c -E "cassandra.* 2/2")
        if [ $cassandra_count -eq $1 ]; then
            kubectl -n workshop exec cassandra-0 -- nodetool status
            cassandra_un=$(kubectl -n workshop exec cassandra-0 -- nodetool status | grep -c -E "^UN")
            if [ $cassandra_un -eq $1 ]; then
                return
            fi
        fi

        sleep 30s
    done
}

print_curr_env() {
    #cassandra_count=$(kubectl -n gs-test get po | grep -c -E "cassandra.* 2/2")
    cassandra_count=$(kubectl -n workshop get po | grep -c -E "cassandra.* 1/1")
    fp_count=$(kubectl -n workshop get po | grep -c -E "fp.* 1/1")
    echo "curr time" $(date '+%Y-%m-%d %H:%M:%S') "cassandra:" $cassandra_count "fp:" $fp_count
}

for ((i = 3; i <= 3; i = i + 3)); do
    # replica_cassandra $i
    # for ((j = 1; j <= i * 2 + 2; j = j + 2)); do
        # content_num=$(($j * 12))
        # echo "begin test expect cassandra:" $i "fp:" $j
        # replica_fp $j
        print_curr_env
        # echo "content" $content_num
        j=$(kubectl -n workshop get po | grep -c -E "fp.* 1/1")
        # 热数据
        echo "开始热数据"
        wrk --latency -s $baseDir/lua/script-zxb.lua "http://10.1.9.188:52352" -t 3 -c $i -d 30s -- $baseDir/data/du.json test false
        # ./wrk --latency -s script-zxb.lua "https://sandbox.dv-api.com" -t 3 -c $content_num -d 30s -- data/test.json GS_STRESS_TEST false
        echo "结束热数据"

        # 正式测试
        echo "begin without update 35fp test expect cassandra:" $i "fp:" $j
        wrk --latency -s $baseDir/lua/script-zxb.lua "http://10.1.9.188:52352" -t 3 -c $i -d 2m -- $baseDir/data/du.json test false

        echo "begin without update 1fp test expect cassandra:" $i "fp:" $j
        wrk --latency -s $baseDir/lua/script-zxb.lua "http://10.1.9.188:52352" -t 3 -c $i -d 2m -- $baseDir/data/du.json test false

        echo "begin with update 1fp test expect cassandra:" $i "fp:" $j
        wrk --latency -s $baseDir/lua/script-zxb.lua "http://10.1.9.188:52352" -t 3 -c $i -d 2m -- $baseDir/data/du.json test true

        echo "begin with update 1fp test expect cassandra:" $i "fp:" $j
        wrk --latency -s $baseDir/lua/script-zxb.lua "http://10.1.9.188:52352" -t 3 -c $i -d 2m -- $baseDir/data/du.json test true
    #done
    alarm "fp:$j,cassandra: $i,此轮测试结束"
done




exit

#print_curr_env
#replica_fp 2
#replica_cassandra 4
#./wrk --latency -s script-zxb.lua "https://sandbox.dv-api.com" -t 3 -c 10 -d 10s -- data/test.json GS_STRESS_TEST false
#./wrk --latency -s script-zxb.lua "https://sandbox.dv-api.com" -t 3 -c 10 -d 10s -- data/test.json GS_STRESS_TEST false
