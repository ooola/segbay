#!/bin/bash

ID=$$
URL='http://127.0.0.1:9999/bids/'

curl -v -H "Accept: application/json" -X GET "$URL" 
