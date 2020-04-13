#!/bin/sh

task=$1
db_schema=$2
count=$3
letter=$4
repcount=$5
threads=$6

#export YCSB_HOME=/driver

if [ "$task" == "" ]; then
    echo ""
    echo "Usage: ycsb_workload_gen.sh load|run [schema] [opsCount] [a|b|c|d|e|f|g] [repCount]"
    echo ""
    echo "    db schema"
    echo "             (default=user)"
    echo "             (to set, must also set DB_NAME and AGENT env vars)"
    echo "             (  otherwise will use nuo.properties values      )"
    echo ""
    echo "    opsCount =  num of rows for load operation"
    echo "             num of transactions for run operation"
    echo "             (default=10,000)"
    echo ""
    echo "    workload a =  50/50 read update"
    echo "             b =  95/5  read update"
    echo "             c =  100   read"
    echo "             d =  95/5  read insert"
    echo "             e =  95/5  scans insert"
    echo "             f =  50/50 read read/modify"
    echo ""
    echo "    repCount = num of times to run the benchmark"
    echo "             (default=1)"
    echo "             (continuous=0)"
    echo ""
    echo "    threads = Number of SQL statements processed per connection"
    echo "             (default=10)"
    echo ""
    echo "    override ENV variables"
    echo "             DB_NAME"
    echo "             DB_USER"
    echo "             DB_PASSWORD"
    echo ""
    exit
fi

if [ "$letter" == "" ]; then
    letter="c"
fi

if [ "$db_schema" == "" ]; then
    db_schema="user"
fi

if [ "$task" == "load" ]; then
    counttype="recordcount"
else
    counttype="operationcount"
fi

if [ "$repcount" == "" ]; then
    repcount=1
elif [ "$repcount" == "0" ]; then
    repcount=9999999
fi

# NTJ: The YCSB app/script parses the command line breaking on spaces - so LB_QUERY cannot work.
# LOADBALANCER_QUERY
# if [ -n "$LB_QUERY" ]; then
#     LOADBALANCER_QUERY="&LB_QUERY=${LB_QUERY}"
# fi

LOADBALANCER_POLICY=""
if [ -n "${LB_POLICY}" ]; then
    LOADBALANCER_POLICY="&LBPolicy=${LB_POLICY}"
fi

#check for nuo.properties file overrides
if [ "$PEER_ADDRESS" != "" ] && [ "$DB_NAME" != "" ]; then
    AGENT_DB_NAME_VALUE="-p db.url=jdbc:com.nuodb://${PEER_ADDRESS}/${DB_NAME}?schema=${db_schema}"

    if [ "$TE_DIRECT" == "true" ]; then
        AGENT_DB_NAME_VALUE="${AGENT_DB_NAME_VALUE}&direct=true"
    else
        AGENT_DB_NAME_VALUE="${AGENT_DB_NAME_VALUE}${LOADBALANCER_QUERY}${LOADBALANCER_POLICY}"
    fi
fi

if [ "$DB_USER" != "" ]; then
    DB_USER_VALUE="-p db.user=${DB_USER}"
fi

if [ "$DB_PASSWORD" != "" ]; then
    DB_PASSWORD_VALUE="-p db.passwd=${DB_PASSWORD}"
fi

rm ycsb_$letter.log > /dev/null 2>&1

for i in $(seq 1 $repcount); do

echo "---------------------------------------------------------" >> ycsb_$letter.log
echo "----- run number $i at "  `date`                               >> ycsb_$letter.log
echo "---------------------------------------------------------" >> ycsb_$letter.log

echo "==============================================================================================="
echo "YCSB_HOME=$YCSB_HOME"

echo "$YCSB_HOME/bin/ycsb.sh $task jdbc \
         -P $YCSB_HOME/nuo.properties \
         -P $YCSB_HOME/workloads/workload$letter \
         -p $counttype=$count \
         $AGENT_DB_NAME_VALUE \
         $DB_USER_VALUE \
         $DB_PASSWORD_VALUE \
         -p threads=2 -s >> ycsb_$letter.log 2>&1
"

$YCSB_HOME/bin/ycsb.sh $task jdbc \
         -P $YCSB_HOME/nuo.properties \
         -P $YCSB_HOME/workloads/workload$letter \
         -p $counttype=$count \
         $AGENT_DB_NAME_VALUE \
         $DB_USER_VALUE \
         $DB_PASSWORD_VALUE \
         -p threads=2 -s >> ycsb_$letter.log 2>&1

done

echo ""
exit 0
