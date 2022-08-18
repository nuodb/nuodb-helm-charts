#!/bin/sh

#This script is performing the steps here https://github.com/dbeaver/cloudbeaver/wiki/Adding-new-database-drivers#adding-drivers-in-cloudbeaver
#It may break with future releases of cloudbeaver, but the benefit is that we don't need to mantain our own version of the cloudbeaver image for use with NuoDB.

#Copy the jars onto the shared volume, so that they can be managed in the next init container without having to install more sofware
cp -r /opt/cloudbeaver/server/plugins/* /opt/cloudbeaver/cloudbeaver-jars/;
echo "jars copied to /opt/cloudbeaver/cloudbeaver-jars/";