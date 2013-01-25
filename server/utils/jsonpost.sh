#!/bin/bash

ID=$$
URL='http://127.0.0.1:9999/bids/'
BODY='{ "Itemid": "'$ID'", "Amount": "100", "Time": "2012-10-13T19:35:18.000Z", "SecondsBefore": "10"}'

curl -v -H "Accept: application/json" -H "Content-type: application/json" -X POST -d "$BODY" "$URL"
