#!/bin/sh

#This script is performing the steps here https://github.com/dbeaver/cloudbeaver/wiki/Adding-new-database-drivers#adding-drivers-in-cloudbeaver
#It may break with future releases of cloudbeaver, but the benefit is that we don't need to mantain our own version of the cloudbeaver image for use with NuoDB.

#Grab any jars that were updated
cp /opt/cloudbeaver/cloudbeaver-jars/* /opt/cloudbeaver/server/plugins/;

#Now run the original cloudbeaver script
./run-server.sh;