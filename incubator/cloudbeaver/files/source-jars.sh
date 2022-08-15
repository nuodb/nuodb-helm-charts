#!/bin/sh

#This script is performing the steps here https://github.com/dbeaver/cloudbeaver/wiki/Adding-new-database-drivers#adding-drivers-in-cloudbeaver
#It may break with future releases of cloudbeaver, but the benefit is that we don't need to mantain our own version of the cloudbeaver image for use with NuoDB.

#Move the two jars that need modifying onto the shared volume, so that they can be managed in the next init container without having to install more sofware
CLOUDBEAVER_BUNDLE_NAME=$(find /opt/cloudbeaver/server/plugins -name "io.cloudbeaver.resources.drivers.base*");
cp $CLOUDBEAVER_BUNDLE_NAME /opt/cloudbeaver/cloudbeaver-jars/;
echo "$CLOUDBEAVER_BUNDLE_NAME copied to /opt/cloudbeaver/cloudbeaver-jars";

JKISS_BUNDLE_NAME=$(find /opt/cloudbeaver/server/plugins -name "org.jkiss.dbeaver.ext.generic*");
cp $JKISS_BUNDLE_NAME /opt/cloudbeaver/cloudbeaver-jars/;
echo "$JKISS_BUNDLE_NAME copied to /opt/cloudbeaver/cloudbeaver-jars";