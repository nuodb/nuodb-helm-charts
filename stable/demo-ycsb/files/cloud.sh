#!/bin/sh

function cloud_provider()
{
    local provider="$(cat /sys/devices/virtual/dmi/id/sys_vendor)"
    if echo "${provider}" | grep -q .*Amazon.* ; then
        echo "aws"
    elif echo "${provider}" | grep -q .*Google.* ; then
        echo "gce"
    elif echo "${provider}" | grep -q .*Microsoft.* ; then
        echo "azure"
    else
        echo "unknown"
    fi
}

# echo $(cloud_provider)

function cloud_zone()
{
    local provider="$(cloud_provider)"
    local zone="unknown"
    if [ "${provider}" == "aws" ] ; then
        zone="$(curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone/)"
    elif [ "${provider}" == "gce" ] ; then
        zone=$(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/zone)
        zone=$(basename ${zone})
    elif [ "${provider}" == "azure" ] ; then
        zone=$(curl -s -H Metadata:true "http://169.254.169.254/metadata/instance/compute/zone?api-version=2018-04-02&format=text")
        [ "${zone}" == "" ] && zone=$(curl -s -H Metadata:true "http://169.254.169.254/metadata/instance/compute/location?api-version=2018-04-02&format=text")
    fi
    echo ${zone}
}

# echo $(cloud_zone)

function cloud_region()
{
    local provider="$(cloud_provider)"
    local zone="$(cloud_zone)"
    local region="unknown"
    if [[ "${provider}" == "azure" ]] ; then
        region=$(curl -s -H Metadata:true "http://169.254.169.254/metadata/instance/compute/location?api-version=2018-04-02&format=text")
    elif [[ "${provider}" == "aws" ]] ; then
        region=$(echo "$zone" | sed s'/.$//' )
    elif [[ "${provider}" == "gce" ]] ; then
        region=$(echo "$zone" | sed s'/..$//' )
    fi
    echo ${region}
}

# echo $(cloud_region)

function cloud_hostname()
{
    local provider="$(cloud_provider)"
    local hostname=$(hostname -f)
    if [[ "${provider}" == "azure" ]] ; then
        hostname=$(curl -s -H Metadata:true "http://169.254.169.254/metadata/instance/compute/name?api-version=2018-04-02&format=text")
    elif [[ "${provider}" == "aws" ]] ; then
        hostname="$(curl -s http://169.254.169.254/latest/meta-data/local-hostname)"
    elif [[ "${provider}" == "gce" ]] ; then
        hostname=$(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/hostname)
    fi
    echo ${hostname}
}

# echo $(cloud_hostname)
