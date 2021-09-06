#!/bin/bash

#while true ; do
#    DATE_TIME_STAMP=$(date --utc +%FT%T.%3NZ)
#    BODY='{"user_id":"1","location_id":"234","checkin_timestamp":"'"${DATE_TIME_STAMP}"'"}'
#    curl -k https://dapr-checkin.kainiindustries.net/v1.0/invoke/frontend/method/checkin -X POST -d $BODY
    # seq 1 10 | xargs -n1 -P10  curl -k https://dapr-checkin.kainiindustries.net/v1.0/invoke/frontend/method/checkin -X POST -d $BODY
#done

function conc () {
    cmd=("${@:3}")
    seq 1 "$1" | xargs -n1 -P"$2" "${cmd[@]}"
}

DATE_TIME_STAMP=$(date --utc +%FT%T.%3NZ)
BODY='{"user_id":"1","location_id":"234","checkin_timestamp":"'"${DATE_TIME_STAMP}"'"}'
conc 10000 10 curl -k https://dapr-checkin.kainiindustries.net/v1.0/invoke/frontend/method/checkin -X POST -d $BODY
