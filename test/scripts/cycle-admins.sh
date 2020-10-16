#!/bin/bash

#
# script to repeatedly cycle the APs,
# verifying that the raft state remains consistent after each iteration

# all errors are fatal
set -e

: ${KUBE_NAMESPACE:=nuodb}

##### wait for an event #####
function wait_for() {
   sleepTime=2

   LOOP_RETRY=$(( ${3:-30}/sleepTime ))

   echo "waiting for $1 == $2..."

   until current=$($1); [ "$current" = "$2" ]; do
      echo "waitfor: found=$current"

      if [ $((LOOP_RETRY--)) == 0 ]; then
         echo "waitfor: timed out on $1 == $2"
         echo "current = $current"

         for ap in $aplist; do
            $kn exec -it $ap -- nuocmd show domain
         done

         [ -n "$DEBUG" ] && sleep 84600
         exit 1
      fi

      sleep $sleepTime;
   done

   echo "waitfor: found $current; countdown=$LOOP_RETRY"
}
#####

# set the kubectl command
kn="kubectl -n $KUBE_NAMESPACE"

# get the admin-0 AP
apzero=$( $kn get pods | grep -Eo '^admin-[^ ]+-0 ' )

# get the list of AP pods
apset=$( $kn get pods | grep -Eo '^admin-[^ ]+' )

# reverse the order of the list
aplist=""
for ap in $apset; do
    aplist="$ap $aplist"
done

get_state="$kn exec -it $apzero -- nuocmd --show-json-fields id,connectedState.state,peerState,peerMemberState get servers"

# get the current state
initial_state=$( $get_state )

startrun=$(date +%s)
endrun=$(( startrun + ${1:-120} ))

# repeat the test for the specified time
while [ `date +%s` -lt $endrun ]; do
    # cycle through the APs
    for ap in $aplist; do
        $kn delete pod $ap

        #wait_for "$kn exec -it $ap -- nuocmd check servers" "" 60
        wait_for "$get_state" "$initial_state" 60

        for ap in $aplist; do
            $kn exec -it $ap -- nuocmd get servers
        done

        $kn exec -it $apzero nuocmd show domain
    done

    sleep 60
done
