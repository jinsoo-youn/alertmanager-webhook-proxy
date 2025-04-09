#!/bin/bash

echo "Sending test alert to alertmanager-webhook-proxy..."

curl -s -X POST http://0.0.0.0:8080/webhook \
     -H "Content-Type: application/json" \
     -d @sample-alert.json \
     | jq .

echo "Done"