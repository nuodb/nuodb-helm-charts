#!/bin/sh

#########################################
#
# expected params
#
#    YCSB_WORKLOAD = work load
#    PEER_ADDRESS = admin service load balancer address
#    DB_NAME
#    DB_USER
#    DB_PASSWORD

set -e

: ${MAX_DELAY:=240000}

# YCSB_THREADS=10
# if [ -n "${NO_OF_PROCESSES}" ]; then
#     YCSB_THREADS="${NO_OF_PROCESSES}"
# fi

# NUMOFROWS=10000
# if [ -n "${NO_OF_ROWS}" ]; then
#     NUMOFROWS="${NO_OF_ROWS}"
# fi

# NUMOFITERATIONS=0
# if [ -n "${NO_OF_ITERATIONS}" ]; then
#     NUMOFITERATIONS="${NO_OF_ITERATIONS}"
# fi

# set default values
: ${NO_OF_PROCESSES:=10}
: ${NO_OF_ROWS:=10000}
: ${NO_OF_ITERATIONS:=10000}
: ${OPS_PER_ITERATION:=10000}

: ${YCSB_WORKLOAD:=b}

# wait to be able to create connections, and timeout after 4 minutes...
java -cp /driver/lib/simple-sql-cli.jar com.nuodb.scripts.SqlDial --broker ${PEER_ADDRESS} --database ${DB_NAME} --user ${DB_USER} --password ${DB_PASSWORD} --timeout ${MAX_DELAY} --connection-property direct=$TE_DIRECT

#check of table is populated
tablecount=$( java -jar /driver/lib/simple-sql-cli.jar --broker ${PEER_ADDRESS} --database ${DB_NAME} -P schema=user1 -P direct=$TE_DIRECT --user ${DB_USER} --password ${DB_PASSWORD} --file /driver/check_usertable.sql | grep 'usertable:' | awk -F ',' '{ print $2; }' )
if [ "$tablecount" == "0" ]; then
    #create table
    /driver/create_usertable.sh 1

    #populate table
    echo "Table ROWS: ${NO_OF_ROWS}"
    /driver/ycsb_gen.sh load user1 ${NO_OF_ROWS}
fi


for i in $(seq 1 $NO_OF_PROCESSES);
do
    /driver/ycsb_gen.sh run user1 ${OPS_PER_ITERATION} ${YCSB_WORKLOAD} ${NO_OF_ITERATIONS} 2 &
done

#/driver/ycsb_gen.sh run user1 ${NUMOFROWS} ${YCSB_WORKLOAD} ${NUMOFITERATIONS} ${YCSB_THREADS} &

while [ true ]; do
    sleep 1200
done