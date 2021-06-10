#!/bin/sh

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
java -jar /driver/lib/simple-sql-cli.jar --broker ${PEER_ADDRESS} --user ${DB_USER} --password ${DB_PASSWORD} \
   --database ${DB_NAME} -P schema=hockey -P direct=$TE_DIRECT --file /driver/create-db.sql

java -jar /driver/lib/simple-sql-cli.jar --broker ${PEER_ADDRESS} --user ${DB_USER} --password ${DB_PASSWORD} \
   --database ${DB_NAME} -P schema=hockey -P direct=$TE_DIRECT --file /driver/Players.sql

java -jar /driver/lib/simple-sql-cli.jar --broker ${PEER_ADDRESS} --user ${DB_USER} --password ${DB_PASSWORD} \
   --database ${DB_NAME} -P schema=hockey -P direct=$TE_DIRECT --file /driver/Scoring.sql

java -jar /driver/lib/simple-sql-cli.jar --broker ${PEER_ADDRESS} --user ${DB_USER} --password ${DB_PASSWORD} \
   --database ${DB_NAME} -P schema=hockey -P direct=$TE_DIRECT --file /driver/Teams.sql

echo "Done"
echo ""


