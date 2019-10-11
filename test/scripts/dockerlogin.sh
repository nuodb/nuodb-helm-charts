#!/usr/bin/env bash

# script to create dockerlogin credentials for docker private registries

: ${SECRET_NAME:="dockerlogin"}
: ${DOCKER_SERVER:="docker.io"}
: ${DOCKER_USERNAME?"Must set username"}
: ${DOCKER_PASSWORD?"Must set password variable"}
: ${DOCKER_EMAIL?"Must set email address"}

DOCKER_CREDS="${DOCKER_USERNAME}:${DOCKER_PASSWORD}"
DOCKER_AUTH=$(echo ${DOCKER_CREDS} | base64)

cat > pullsecret.yaml <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: ${SECRET_NAME}
type: kubernetes.io/dockerconfigjson
stringData:
  .dockerconfigjson: |-
    {
       "auths" : {
         "${DOCKER_SERVER}" : {
           "username" : "${DOCKER_USERNAME}",
           "password" : "${DOCKER_PASSWORD}",
           "email"    : "${DOCKER_EMAIL}",
           "auth"     : "${DOCKER_AUTH}"
         }
       }
    }
EOF

kubectl delete -f pullsecret.yaml
kubectl create -f pullsecret.yaml
