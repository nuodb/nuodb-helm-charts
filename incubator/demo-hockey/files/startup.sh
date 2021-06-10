#!/bin/sh

#    Expected Parameters
#
#    PEER_ADDRESS = admin service load balancer address
#    DB_NAME
#    DB_USER
#    DB_PASSWORD

set -e

# set default values
: ${MAX_DELAY:=240000}


# wait to be able to create connections, and timeout after 4 minutes...
java -cp /driver/lib/simple-sql-cli.jar com.nuodb.scripts.SqlDial \
   --broker ${PEER_ADDRESS} --database ${DB_NAME} --user ${DB_USER} --password ${DB_PASSWORD} \
   --timeout ${MAX_DELAY} --connection-property direct=$TE_DIRECT

#check of table is populated
tablecount=$( java -jar /driver/lib/simple-sql-cli.jar --broker ${PEER_ADDRESS} \
   --database ${DB_NAME} -P schema=hockey -P direct=$TE_DIRECT --user ${DB_USER} --password ${DB_PASSWORD} \
   --file /driver/check-players-table.sql | grep 'players:' | awk -F ',' '{ print $2; }' )
if [ "$tablecount" == "0" ]; then
    #create hockey schema
    /driver/create-hockey-schema.sh > /dev/null
fi

# For Docker Compose, Run, or Kubernetes Job, exit the container if only creating the schema
if [ "$CREATE_SCHEMA_ONLY" == "true" ]; then
   exit 0
fi

# run workload

while [ true ]; do
   java -jar /driver/lib/simple-sql-cli.jar --broker ${PEER_ADDRESS} --user ${DB_USER} --password ${DB_PASSWORD} \
     --database ${DB_NAME} -P schema=hockey -P direct=$TE_DIRECT --file /driver/hockey-workload.sql
done

exit

# run workload using SimpleDriver                                 
                              
java -jar /driver/lib/SimpleDriver.jar -url jdbc:com.nuodb://${PEER_ADDRESS}/${DB_NAME} \              
  -user ${DB_USER} -password ${DB_PASSWORD} \                                                          
  -property schema:hockey -property direct:true \                                                                
  -time 0 -threads 0 -minthreads 5 \                                                                                                      
  -sql file:/driver/hockey-workload.sql                                                                                           

