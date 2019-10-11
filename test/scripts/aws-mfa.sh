#!/usr/bin/env bash

: ${MFA_URN:="arn:aws:iam::650436259446:mfa/demo-operator"}

token=

export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=
export AWS_SESSION_TOKEN=

function usage {
  echo "usage: $0 [[-t token] | [-h]]"
}

function set-mfa {
  json=`aws --output json sts get-session-token --serial-number ${MFA_URN} --token-code ${token}`
  echo ${json} > ~/.aws/mfa.json
  access_key_id=`cat ~/.aws/mfa.json | jq -e '.Credentials.AccessKeyId' | tr -d '"'`
  secret_access_key=`cat ~/.aws/mfa.json | jq -e '.Credentials.SecretAccessKey' | tr -d '"'`
  session_token=`cat ~/.aws/mfa.json | jq -e '.Credentials.SessionToken' | tr -d '"'`
cat << EOF
export AWS_ACCESS_KEY_ID=${access_key_id}
export AWS_SECRET_ACCESS_KEY=${secret_access_key}
export AWS_SESSION_TOKEN=${session_token}
EOF
}

while [ "$1" != "" ]; do
  case $1 in
    -t | --token )
      shift
      token=$1
      set-mfa
      ;;
    -h | --help )
      usage
      exit
      ;;
    * )
      usage
      exit 1
  esac
  shift
done
