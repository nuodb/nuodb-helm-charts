#!/bin/sh

#This script is performing the steps here https://github.com/dbeaver/cloudbeaver/wiki/Adding-new-database-drivers#adding-drivers-in-cloudbeaver
#It may break with future releases of cloudbeaver, but the benefit is that we don't need to mantain our own version of the cloudbeaver image for use with NuoDB.

#Get the driver filename to fetch, or default to the current version at time of writing
NUODB_DRIVER_VERSION="${NUODB_DRIVER_VERSION:=nuodb-jdbc-23.0.0.jar}";
NUODB_DRIVER_NAME="nuodb-jdbc-$NUODB_DRIVER_VERSION.jar";

rm -rf /opt/cloudbeaver/drivers/nuodb/*;
if [ "$DOWNLOAD_DRIVER" = "true" ]
then
  #Download the driver from maven and place it in the cloudbeaver driver directory
  echo "Downloading NuoDB Driver $NUODB_DRIVER_NAME from maven";
  wget -O /opt/cloudbeaver/drivers/nuodb/$NUODB_DRIVER_NAME https://repo1.maven.org/maven2/com/nuodb/jdbc/nuodb-jdbc/$NUODB_DRIVER_VERSION/$NUODB_DRIVER_NAME;
else
  #Decode it from the configmap
  echo "Copying NuoDB Driver $NUODB_DRIVER_NAME from configmap";
  cp /opt/cloudbeaver/config/nuodb/$NUODB_DRIVER_NAME /opt/cloudbeaver/drivers/nuodb/$NUODB_DRIVER_NAME;
fi;

#Replace the two plugin.xml files containing the additional NuoDB config
CLOUDBEAVER_BUNDLE_NAME=$(find /opt/cloudbeaver/cloudbeaver-jars/ -name "io.cloudbeaver.resources.drivers.base*");
jar -uf $CLOUDBEAVER_BUNDLE_NAME -C /opt/cloudbeaver/config/nuodb/io.cloudbeaver.resources.drivers.base plugin.xml;
echo "$CLOUDBEAVER_BUNDLE_NAME updated.";

JKISS_BUNDLE_NAME=$(find /opt/cloudbeaver/cloudbeaver-jars/ -name "org.jkiss.dbeaver.ext.generic*");
jar -uf $JKISS_BUNDLE_NAME -C /opt/cloudbeaver/config/nuodb/org.jkiss.dbeaver.ext.generic plugin.xml;
echo "$JKISS_BUNDLE_NAME updated.";


