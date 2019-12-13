#!/bin/sh


YCSB_SCHEMA_COUNT=$1

if [ "$YCSB_SCHEMA_COUNT" == "" ]; then
  echo "usage: create_usertable.sh numberOfSchemas"
  exit 0
fi

if [ "$DB_NAME" == "" ]; then
  DB_NAME="test"
fi

if [ "$DB_USER" == "" ]; then
  DB_USER="dba"
fi

if [ "$DB_PASSWORD" == "" ]; then
  DB_PASSWORD="goalie"
fi

if [ "$PEER_ADDRESS" == "" ]; then
  AGENT="nuoadmin1"
fi

echo ""
for i in $(seq 1 $YCSB_SCHEMA_COUNT); do
    echo "Creating ycsb table in schema user${i}"
    #java -cp /driver/jar/ nuodb.SimpleDriver -config /driver/ycsb.props -p DB_SCHEMA=user${i} >> /var/tmp/create_usertable.log 2>&1
    java -jar /driver/lib/simple-sql-cli.jar --broker ${PEER_ADDRESS} --user ${DB_USER} --password ${DB_PASSWORD} --database ${DB_NAME} -P schema=user${i} -P direct=$TE_DIRECT --file /driver/create_usertable.sql
done

echo "Done"
echo ""